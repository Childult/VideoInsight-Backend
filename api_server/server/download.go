package server

import (
	"fmt"
	"path/filepath"
	"strconv"
	"swc/mongodb"
	"swc/mongodb/job"
	"swc/mongodb/resource"
	"swc/util"
	"time"
)

func creatResource(job *job.Job) {
	// 构建资源
	resource := resource.Resource{
		URL:      job.URL,
		Status:   util.ResourceDownloading,
		Location: filepath.Join(util.SavePath, strconv.FormatInt(time.Now().Unix(), 10)),
	}

	// 检查资源是否存在
	exists := mongodb.HaveExisted(resource)
	if exists {
		job.SetStatus(util.JobExisted)
		go JobSchedule(job)
		return
	}
	// 首次写入数据库
	mongodb.InsertOne(resource)
	job.SetStatus(util.JobDownloadMedia)
	go JobSchedule(job)
}

func mediaDownload(job *job.Job) {
	// 获取资源信息
	resource, err := resource.GetByKey(job.URL)
	if err != nil {
		// 获取资源出错
		job.SetStatus(util.JobErrFailedToFindResource)
		go JobSchedule(job)
		return
	}

	// 构建视频下载对象
	python := PyWorker{
		PackagePath: filepath.Join(util.WorkSpace, "video_getter"),
		FileName:    "main",
		MethodName:  "download_video",
		Args: []string{
			SetArg(resource.URL),
			SetArg(resource.Location),
		},
	}

	go python.Call(job, downloadHandle)
}

func downloadHandle(job *job.Job, result []string) {
	// 获取资源信息
	r, err := resource.GetByKey(job.URL)
	if err != nil {
		// 获取资源出错
		job.SetStatus(util.JobErrFailedToFindResource)
		go JobSchedule(job)
		return
	}

	// 是否下载成功
	if len(result) != 1 {
		// 下载失败
		r.SetStatus(util.ResourceErrDownloadFailed)
		job.SetStatus(util.JobErrDownloadFailed)
		go JobSchedule(job)
		return
	}

	// 下载成功, 更新状态
	r.VideoPath = result[0]
	r.SetStatus(util.ResourceExtracting)
	job.SetStatus(util.JobExtractAudio)
	go JobSchedule(job)
}

func extractAudio(job *job.Job) {
	// 获取资源信息
	r, err := resource.GetByKey(job.URL)
	if err != nil {
		// 获取资源出错
		job.SetStatus(util.JobErrFailedToFindResource)
		go JobSchedule(job)
		return
	}

	// 构建音频提取对象
	python := PyWorker{
		PackagePath: filepath.Join(util.WorkSpace, "video_analysis"),
		FileName:    "extract_audio",
		MethodName:  "extract_audio",
		Args: []string{
			SetArg(filepath.Join(r.Location, r.VideoPath)),
		},
	}

	// 提取音频
	go python.Call(job, extractHandle)
}

func extractHandle(job *job.Job, result []string) {
	// 获取资源信息
	r, err := resource.GetByKey(job.URL)
	if err != nil {
		// 获取资源出错
		job.SetStatus(util.JobErrFailedToFindResource)
		go JobSchedule(job)
		return
	}

	// 是否成功提取音频
	if len(result) != 1 {
		// 提取失败
		r.SetStatus(util.ResourceErrExtractFailed)
		job.SetStatus(util.JobErrExtractFailed)
		go JobSchedule(job)
		return
	}

	// 音频提取成功, 更新状态
	r.AudioPath = result[0]
	r.SetStatus(util.ResourceCompleted)
	job.SetStatus(util.JobExtractDone)
	go JobSchedule(job)
}

func waitDownload(job *job.Job) {
	// 获取资源信息
	fmt.Println("========================= 获取资源信息 ===================================")
	r, err := resource.GetByKey(job.URL)
	if err != nil {
		// 获取资源出错
		job.SetStatus(util.JobErrFailedToFindResource)
		go JobSchedule(job)
		return
	}

	for {
		if r.Status == util.ResourceCompleted {
			job.SetStatus(util.JobExtractDone)
			go JobSchedule(job)
		} else if r.Status > util.ResourceCompleted {
			fmt.Println("========================= 资源下载出错 ===================================")

		} else {

			fmt.Println("========================= 资源已经存在, 等待下载完成 ===================================")
			time.Sleep(time.Second * 5)
		}
	}
}

func extractAbstract(job *job.Job) {
	// 进行文本分析
	go textAnalysis(job)
	// 进行视频分析
	// go videoAnalysis(job)
}
