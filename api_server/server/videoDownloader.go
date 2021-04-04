package server

import (
	"path/filepath"
	"strconv"
	"strings"
	"swc/data/resource"
	"swc/data/task"
	"swc/dbs/mongodb"
	"swc/dbs/redis"
	"swc/logger"
	"swc/util"
	"sync"
	"time"
)

var creatMu sync.Mutex

func creatResource(task *task.Task) (r *resource.Resource) {
	logger.Info.Printf("[创建资源] 开始: %+v.\n", task)
	// 构建资源
	r = &resource.Resource{
		URL:      task.URL,                                                                     // 资源链接
		Status:   util.ResourceCreated,                                                         // 状态: 下载中
		Location: filepath.Join(util.Location, strconv.FormatInt(time.Now().Unix(), 10)) + "/", // 保存地址, 以时间戳为文件夹
	}

	// 检查资源是否存在
	// 加锁, 保证一份资源不会被多次创建
	creatMu.Lock()
	if redis.Exists(r) {
		// 存在的资源不重复下载
		creatMu.Unlock()
		logger.Warning.Println("[创建资源] 资源已存在")
		task.Status = util.JobResourceExisted
		redis.FindOne(r)
	} else if mongodb.Exists(r) {
		// 存在的资源不重复下载
		creatMu.Unlock()
		logger.Warning.Println("[创建资源] 资源已存在")
		task.Status = util.JobResourceExisted
		mongodb.FindOne(r)
	} else {
		// 不存在则创建资源
		redis.InsertOne(r)
		creatMu.Unlock()
		logger.Info.Printf("[创建资源] 资源创建成功:%+v.\n", r)
		task.Status = util.JobToDownloadMedia
	}
	redis.UpdataOne(task)
	return
}

func mediaDownload(task *task.Task, r *resource.Resource) {
	logger.Info.Printf("[视频下载] 开始: %+v.\n", task)
	// 开始下载
	r.Status = util.ResourceDownloading
	redis.UpdataOne(r)

	// 构建视频下载对象
	python := PyWorker{
		PackagePath: filepath.Join(util.WorkSpace, "video_getter"), // python 包地址
		FileName:    "api",                                         // 文件名
		MethodName:  "download_video",                              // 调用函数
		Args: []string{ // 实参
			SetArg(r.URL),      // 资源链接
			SetArg(r.Location), // 保存路径
		},
	}

	result := python.Call()

	// 是否下载成功
	if len(result) == 0 {
		logger.Error.Println("[视频下载] 下载失败.")
		r.Status = util.ResourceErrDownloadFailed
		task.Status = util.JobErrDownloadFailed
	} else {
		// 下载成功, 更新状态
		pythonReturn := strings.Join(result, "")
		logger.Debug.Printf("[视频下载] 视频下载成功: %+v.\n", pythonReturn)
		r.VideoPath = pythonReturn

		r.Status = util.ResourceDownloadDone
		task.Status = util.JobToExtractAudio
	}

	redis.UpdataOne(r)
	redis.UpdataOne(task)
}
