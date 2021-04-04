package resource_builder

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"swc/data/resource"
	"swc/dbs/mongodb"
	"swc/dbs/redis"
	"swc/logger"
	"swc/server/python"
	"swc/util"
	"sync"
	"time"
)

// RSArgs 资源获取时需要输入的参数
type RSArgs struct {
	url  string     // 要获取的资源链接
	back chan error // 用于通知结果的管道
}

// ResourceScheduler 资源获取的调度器
type ResourceScheduler struct {
	mu  sync.Mutex              // 获取资源时加锁, 保证不会重复下载资源
	m   map[string][]chan error // 当资源正在获取时, 有同样的资源进来, 就先保存起来, 以 URL 为键
	chs chan RSArgs             // 调度时需要的参数
}

var rss ResourceScheduler
var onceRSS sync.Once

// RequestResource 请求资源
// url: 想要获取的资源链接
func RequestResource(url string) error {
	// 调度器只会启动一次
	onceRSS.Do(func() {
		rss.m = make(map[string][]chan error)
		rss.chs = make(chan RSArgs)
		go rss.scheduler(rss.chs)
	})

	// 构建参数
	back := make(chan error)
	ch := RSArgs{url: url, back: back}

	// 把参数发送给调度器
	rss.chs <- ch

	// 等待结果
	return <-ch.back
}

// scheduler 进行资源下载的调度, 串行执行
func (ts *ResourceScheduler) scheduler(chs chan RSArgs) {
	// 等待任务
	for ch := range chs {
		// 构建资源对象
		url := ch.url
		r := &resource.Resource{URL: url} // 构建资源

		// 查看是否有已经完成或正在进行的任务
		// 加锁, 保证对管道(m)的访问时串行的
		ts.mu.Lock()
		if redis.Exists(r) {
			// 如果已存在, 取回资源
			redis.FindOne(r)

			// 判断资源状态
			if r.Status == util.ResourceCompleted {
				// 资源已经下载完成, 返回
				ts.mu.Unlock()
				ch.back <- nil
			} else if r.Status > util.ResourceCompleted {
				// 如果资源获取出错了, 返回错误, 用户可以自主选择是否删除重试
				ts.mu.Unlock()
				ch.back <- fmt.Errorf("资源获取失败")
			} else {
				// 已经存在, 但还没有下载完成, 把反馈的管道存起来, 等任务完成时集体通知
				ts.m[url] = append(ts.m[url], ch.back)
				ts.mu.Unlock()
			}
		} else if mongodb.Exists(r) {
			// 只有当任务完成时, 才会持久化到 mongodb
			ts.mu.Unlock()
			ch.back <- nil
		} else {
			// 资源不存在, 则保存管道, 开始执行任务. 先加入 redis 再解锁
			ts.m[url] = append(ts.m[url], ch.back)
			// 更新状态
			r.Status = util.ResourceCreated
			// 保存地址, 以时间戳为文件夹
			r.Location = filepath.Join(util.Location, strconv.FormatInt(time.Now().Unix(), 10)) + "/"
			// 插入数据库, 这样后面相同的链接就不会重复下载了
			redis.InsertOne(r)
			ts.mu.Unlock()

			// 开启下载任务
			go ts.Downloader(r)
		}
	}
}

func (rs *ResourceScheduler) Downloader(r *resource.Resource) {
	// 开始下载视频
	r.Status = util.ResourceDownloading
	redis.UpdataOne(r)

	// 构建视频下载对象
	videoDownloader := python.PyWorker{
		PackagePath: filepath.Join(util.WorkSpace, "video_getter"), // python 包地址
		FileName:    "api",                                         // 文件名
		MethodName:  "download_video",                              // 调用函数
		Args: []string{ // 实参
			python.SetArg(r.URL),      // 资源链接
			python.SetArg(r.Location), // 保存路径
		},
	}

	// 下载
	resultVD := videoDownloader.Call()

	// 是否下载成功
	if len(resultVD) == 0 { // 资源下载失败
		rs.errHappen(r, util.ResourceErrDownloadFailed, "资源下载出错")
		return
	}

	// 下载成功, 更新状态
	vdReturn := strings.Join(resultVD, "")
	logger.Debug.Printf("[视频下载] 视频下载成功: %+v.\n", vdReturn)
	r.VideoPath = vdReturn
	r.Status = util.ResourceDownloadDone
	redis.UpdataOne(r)

	// 开始提取音频
	r.Status = util.ResourceAudioExtracting
	redis.UpdataOne(r)

	// 构建音频提取对象
	audioExtractor := python.PyWorker{
		PackagePath: filepath.Join(util.WorkSpace, "audio_analysis"), // 包名
		FileName:    "api",                                           // 文件名
		MethodName:  "extract_audio",                                 // 调用函数
		Args: []string{ // 实参
			python.SetArg(filepath.Join(r.Location, r.VideoPath)), // 传入视频
		},
	}

	// 提取音频
	resultAE := audioExtractor.Call()
	if len(resultAE) == 0 {
		rs.errHappen(r, util.ResourceErrAudioExtractFailed, "音频提取失败")
		return
	}

	// 音频提取成功, 更新状态
	aeReturn := strings.Join(resultAE, "")
	r.AudioPath = aeReturn
	r.Status = util.ResourceCompleted
	rs.mu.Lock()
	defer rs.mu.Unlock()
	redis.UpdataOne(r)
	mongodb.InsertOne(r)
	for _, back := range rs.m[r.URL] {
		back <- nil
	}
	delete(rs.m, r.URL)
}

func (rs *ResourceScheduler) errHappen(r *resource.Resource, status int32, format string, v ...interface{}) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	err := fmt.Errorf(format, v...)
	for _, back := range rs.m[r.URL] {
		back <- err
	}
	delete(rs.m, r.URL)
	r.Status = status
	redis.UpdataOne(r)
	// 打印错误日志
	logger.Error.Println(err)
	fmt.Println("测试", err)
}
