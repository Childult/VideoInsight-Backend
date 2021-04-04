package absvideo_builder

import (
	"context"
	"fmt"
	"swc/data/absvideo"
	"swc/dbs/mongodb"
	"swc/dbs/redis"
	"swc/logger"
	pb "swc/server/network"
	"swc/util"
	"sync"

	"google.golang.org/grpc"
)

// VAArgs 视频分析需要输入的参数
type VAArgs struct {
	url      string     // 要获取的资源链接
	path     string     // 资源存储的地址
	location string     // 结果存放的地址
	back     chan error // 用于通知结果的管道
}

// VAScheduler 文本分析的调度器结构
type VAScheduler struct {
	mu  sync.Mutex              // 创建任务时加锁, 保证任务唯一
	m   map[string][]chan error // 当一个任务正在进行时, 有同样的任务进来, 就先保存起来, 以 url 为键
	chs chan VAArgs             // 调度参数
}

var vas VAScheduler
var onceVAS sync.Once

// RequestVideoAnalysis 请求对资源进行文本分析
// url: 唯一确定一份视频分析, 将会保存在 absvideo.AbsVideo 中
// path: 资源存储的位置
// location: 结果存放的地址
func RequestVideoAnalysis(url string, path string, location string) error {
	// 调度器只会启动一次
	onceVAS.Do(func() {
		vas.m = make(map[string][]chan error)
		vas.chs = make(chan VAArgs)
		go vas.scheduler(vas.chs)
	})

	// 构建参数
	back := make(chan error)
	ch := VAArgs{url: url, path: path, location: location, back: back}

	// 把任务发送给调度器
	vas.chs <- ch

	// 等待结果
	return <-ch.back
}

// scheduler 视频分析的调度器
func (va *VAScheduler) scheduler(chs chan VAArgs) {
	// 等待任务
	for ch := range chs {
		// 构建文本摘要对象
		url := ch.url
		av := &absvideo.AbsVideo{URL: url}

		// 查看是否有已经完成或正在进行的任务
		// 加锁, 保证对管道(m)的访问时串行的
		va.mu.Lock()
		// 先从 redis 中找, redis 中的数据会过期
		if redis.Exists(av) {
			// 如果已经存在, 取回视频摘要
			redis.InsertOne(av)

			// 判断视频摘要状态
			if av.Status == util.AbsTextComplete {
				// 已经成功, 直接返回
				va.mu.Unlock()
				ch.back <- nil
			} else if av.Status > util.AbsTextComplete {
				// 如果视频分析失败, 返回错误, 用户可以自主选择是否删除重试
				va.mu.Unlock()
				ch.back <- fmt.Errorf("文本分析失败")
			} else {
				// 任务正在进行中, 把反馈的管道存起来, 等任务完成时集体通知
				va.m[url] = append(va.m[url], ch.back)
				va.mu.Unlock()
			}
		} else if mongodb.Exists(av) {
			// 只有当任务完成时, 才会持久化到 mongodb
			va.mu.Unlock()
			ch.back <- nil
		} else {
			// 如果不存在, 则保存管道, 开始执行任务. 先加入 redis 再解锁
			va.m[url] = append(va.m[url], ch.back)
			redis.InsertOne(av)
			va.mu.Unlock()
			go va.videoAnalysis(av, ch.path, ch.location)
		}
	}
}

// videoAbstract 用于存储视频摘要的结果
type videoAbstract struct {
	VAbstract []string `json:"VAbstract"`
	Error     string   `json:"Error"`
}

var tempJobID = "123456"

// videoAnalysis 视频分析
func (va *VAScheduler) videoAnalysis(av *absvideo.AbsVideo, path string, location string) {
	// 设置 gRPC 参数
	address := util.GRPCAddress // gRPC 地址
	jobID := tempJobID          // 任务id, 这里先随便写一个
	videoFile := path           // 视频文件路径
	savaPath := location        // 结果存储路径

	// 连接 gRPC 服务器
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		va.errHappen(av, util.AbsVideoErrGRPCConnectFailed, "gRPC 服务器连接失败: %+v", err)
		return
	}
	defer conn.Close()                     // 关闭连接
	rpc := pb.NewVideoAnalysisClient(conn) // gRPC 调用的句柄

	// 开始调用服务
	rpcResult, err := rpc.GetStaticVideoAbstract(context.TODO(), &pb.VideoInfo{JobId: jobID, File: videoFile, SaveDir: savaPath})
	if err != nil {
		va.errHappen(av, util.AbsVideoErrGRPCallFailed, "gRPC 调用失败: %+v.\n", err)
		return
	}

	// 结果校验
	jobid := rpcResult.GetJobID()
	if jobid != tempJobID {
		va.errHappen(av, util.AbsVideoErrGRPCallJobIDNotMatch, "JobID 不匹配")
		return
	}

	// 提取结果
	var videoPath videoAbstract
	videoPath.VAbstract = rpcResult.GetPicName()
	videoPath.Error = rpcResult.GetError()
	if videoPath.Error != "" {
		va.errHappen(av, util.AbsVideoErrGRPCallFailed, "视频分析失败, 返回错误: <%s>", videoPath.Error)
		return
	}

	// 视频分析成功, 初始化需要存储的数据, 存入数据库
	av.Abstract = videoPath.VAbstract
	av.Status = util.AbsVideoComplete
	va.mu.Lock()
	defer va.mu.Unlock()
	redis.InsertOne(av)
	mongodb.InsertOne(av)
	for _, b := range va.m[av.URL] {
		b <- nil
	}
	delete(va.m, av.URL)
}

func (va *VAScheduler) errHappen(av *absvideo.AbsVideo, status int32, format string, v ...interface{}) {
	va.mu.Lock()
	defer va.mu.Unlock()
	err := fmt.Errorf(format, v...)
	for _, b := range va.m[av.URL] {
		b <- err
	}
	delete(va.m, av.URL)
	av.Status = status
	redis.UpdataOne(av)
	// 打印错误日志
	logger.Error.Println(err)
}
