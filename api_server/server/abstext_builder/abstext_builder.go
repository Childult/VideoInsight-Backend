package abstext_builder

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"swc/data/abstext"
	"swc/dbs/mongodb"
	"swc/dbs/redis"
	"swc/logger"
	"swc/server/python"
	"swc/util"
	"sync"
)

// TAArgs 文本分析需要输入的参数
type TAArgs struct {
	url      string     // 要获取的资源链接
	keyWords []string   // 关键词
	path     string     // 资源存储的地址
	back     chan error // 用于通知结果的管道
}

// TAScheduler 文本分析的调度器结构
type TAScheduler struct {
	mu  sync.Mutex              // 创建任务时加锁, 保证任务唯一
	m   map[string][]chan error // 当一个任务正在进行时, 有同样的任务进来, 就先保存起来, 以文本摘要的哈希为键
	chs chan *TAArgs            // 调度参数
}

var tas TAScheduler
var onceTAS sync.Once

// RequestTextAnalysis 请求对资源进行文本分析
// url, keyWords: 唯一确定一份文本分析, 将会保存在 abstext.AbsText 中
// path: 资源存储的位置
func RequestTextAnalysis(url string, keyWords []string, path string) error {
	logger.Debug.Println("[文本分析]] 收到任务", url, keyWords, path)
	// 调度器只会启动一次
	onceTAS.Do(func() {
		tas.m = make(map[string][]chan error)
		tas.chs = make(chan *TAArgs)
		go tas.scheduler(tas.chs)
	})

	// 构建参数
	back := make(chan error)
	ch := &TAArgs{url: url, keyWords: keyWords, path: path, back: back}

	// 把任务发送给调度器
	tas.chs <- ch

	// 等待结果
	return <-ch.back
}

// scheduler 文本分析的调度器
func (ta *TAScheduler) scheduler(chs chan *TAArgs) {
	// 等待任务
	for ch := range chs {
		// 构建文本摘要对象
		at := abstext.NewAbsText(ch.url, ch.keyWords)

		// 查看是否有已经完成或正在进行的任务
		// 加锁, 保证对管道(m)的访问时串行的
		ta.mu.Lock()
		// 先从 redis 中找, redis 中的数据会过期
		if redis.Exists(at) {
			// 如果已经存在, 取回文本摘要
			redis.InsertOne(at)

			// 判断文本摘要状态
			if at.Status == util.AbsTextComplete {
				// 已经成功, 直接返回
				ta.mu.Unlock()
				ch.back <- nil
			} else if at.Status > util.AbsTextComplete {
				// 如果文本分析失败, 返回错误, 用户可以自主选择是否删除重试
				ta.mu.Unlock()
				ch.back <- fmt.Errorf("文本分析失败")
			} else {
				// 任务正在进行中, 把反馈的管道存起来, 等任务完成时集体通知
				ta.m[at.Hash] = append(ta.m[at.Hash], ch.back)
				ta.mu.Unlock()
			}
		} else if mongodb.Exists(at) {
			// 只有当任务完成时, 才会持久化到 mongodb
			ta.mu.Unlock()
			ch.back <- nil
		} else {
			// 如果不存在, 则保存管道, 开始执行任务. 先加入 redis 再解锁
			ta.m[at.Hash] = append(ta.m[at.Hash], ch.back)
			redis.InsertOne(at)
			ta.mu.Unlock()
			go ta.textAnalysis(at, ch.path, ch.back)
		}
	}
}

// textAbstract 用于存储文本分析的结果, 即文本摘要
type textAbstract struct {
	AText     string   `json:"AText"`
	TAbstract []string `json:"TAbstract"`
	Error     string   `json:"Error"`
}

// textAnalysis 文本分析
func (ta *TAScheduler) textAnalysis(at *abstext.AbsText, path string, back chan error) {
	// 构建文本分析调用对象
	python := python.PyWorker{
		PackagePath: filepath.Join(util.WorkSpace, "text_analysis"), // python 文件所在的包
		FileName:    "api",                                          // 文件名
		MethodName:  "generate_abstract_from_audio",                 // 调用函数名
		Args: []string{ // 调用实参
			python.SetArg(path),
		},
	}

	// 调用文本分析
	result := python.Call()
	if len(result) == 0 { // 未找到结果
		ta.errHappen(at, util.AbsTextErrTextAnalysisReadJSONFailed, "文本分析失败, 无返回结果, python 调用出错")
		return
	}

	pythonReturn := strings.Join(result, "")

	// 找到结果, 从返回结果中提取数据
	var text textAbstract
	err := json.Unmarshal([]byte(pythonReturn), &text)
	if err != nil {
		ta.errHappen(at, util.AbsTextErrTextAnalysisReadJSONFailed, "文本分析失败, JSON 解析错误")
		return
	} else if text.Error != "" {
		ta.errHappen(at, util.AbsTextErrTextAnalysisFailed, "文本分析失败, 文本分析返回错误: <%s>", text.Error)
		return
	}

	// 文本分析成功, 初始化需要存储的数据, 存入数据库
	at.Text = text.AText
	at.Abstract = text.TAbstract
	at.Status = util.AbsTextComplete
	ta.mu.Lock()
	defer ta.mu.Unlock()
	redis.UpdataOne(at)
	mongodb.InsertOne(at)
	for _, b := range ta.m[at.Hash] {
		b <- nil
	}
	delete(ta.m, at.Hash)
}

// errHappen 文本分析过程中出现错误
func (ta *TAScheduler) errHappen(at *abstext.AbsText, status int32, format string, v ...interface{}) {
	ta.mu.Lock()
	defer ta.mu.Unlock()
	err := fmt.Errorf(format, v...)
	for _, b := range ta.m[at.Hash] {
		b <- err
	}
	delete(ta.m, at.Hash)
	at.Status = status
	redis.UpdataOne(at)
	// 打印错误日志
	logger.Error.Println(err)
}
