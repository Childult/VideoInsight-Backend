package server

import (
	"path/filepath"
	"strconv"
	"strings"
	"swc/data/job"
	"swc/data/resource"
	"swc/logger"
	"swc/util"
	"time"
)

func creatResource(job *job.Job) {
	logger.Info.Printf("[创建资源] 开始: %+v.\n", job)
	// 构建资源
	r := resource.Resource{
		URL:      job.URL,                                                                      // 资源链接
		Status:   util.ResourceDownloading,                                                     // 状态: 下载中
		Location: filepath.Join(util.Location, strconv.FormatInt(time.Now().Unix(), 10)) + "/", // 保存地址, 以时间戳为文件夹
	}

	// 检查资源是否存在
	exists := r.ExistInMongodb()
	if exists {
		logger.Warning.Println("[创建资源] 资源已存在")
		job.Status = util.JobExisted
		job.Save()
		go JobSchedule(job)
		return
	}

	// 首次写入数据库
	logger.Debug.Printf("[创建资源] 资源创建成功:%+v.\n", r)
	r.Save()
	job.Status = util.JobDownloadMedia
	job.Save()
	go JobSchedule(job)
}

func mediaDownload(job *job.Job) {
	logger.Info.Printf("[下载视频] 开始: %+v.\n", job)
	// 获取资源信息
	r := resource.Resource{URL: job.URL}
	err := r.Retrieve()
	if err != nil {
		logger.Error.Printf("[下载视频] 获取资源出错: %+v.\n", err)
		job.Status = util.JobErrFailedToFindResource
		job.Save()
		go JobSchedule(job)
		return
	}

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

	go python.Call(job, downloadHandle)
}

func downloadHandle(job *job.Job, result []string) {
	logger.Info.Printf("[下载视频回调] 开始: %+v.\n", job)
	// 获取资源信息
	r := resource.Resource{URL: job.URL}
	err := r.Retrieve()
	if err != nil {
		logger.Error.Printf("[下载视频回调] 获取资源出错: %+v.\n", err)
		job.Status = util.JobErrFailedToFindResource
		job.Save()
		go JobSchedule(job)
		return
	}

	// 是否下载成功
	if len(result) == 0 {
		logger.Error.Println("[下载视频回调] 下载失败.")
		r.Status = util.ResourceErrDownloadFailed
		r.Save()
		job.Status = util.JobErrDownloadFailed
		job.Save()
		go JobSchedule(job)
		return
	}
	pythonReturn := strings.Join(result, "")

	// 下载成功, 更新状态
	logger.Debug.Printf("[下载视频回调] 视频下载成功: %+v.\n", pythonReturn)
	r.VideoPath = pythonReturn
	r.Status = util.ResourceExtracting
	r.Save()
	job.Status = util.JobExtractAudio
	job.Save()
	go JobSchedule(job)
}
