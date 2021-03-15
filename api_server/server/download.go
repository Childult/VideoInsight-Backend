package server

import (
	"path/filepath"
	"strconv"
	"swc/logger"
	"swc/mongodb"
	"swc/mongodb/abstext"
	"swc/mongodb/absvideo"
	"swc/mongodb/job"
	"swc/mongodb/resource"
	"swc/util"
	"time"
)

func creatResource(job *job.Job) {
	logger.Info.Printf("创建资源. [URL: %s] [JobID: %s] [Status: %d]\n", job.URL, job.JobID, job.Status)
	// 构建资源
	resource := resource.Resource{
		URL:      job.URL,
		Status:   util.ResourceDownloading,
		Location: filepath.Join(util.Location, strconv.FormatInt(time.Now().Unix(), 10)) + "/",
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
	return
}

func mediaDownload(job *job.Job) {
	logger.Info.Printf("下载视频. [URL: %s] [JobID: %s] [Status: %d]\n", job.URL, job.JobID, job.Status)
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
		FileName:    "api",
		MethodName:  "download_video",
		Args: []string{
			SetArg(resource.URL),
			SetArg(resource.Location),
		},
	}

	go python.Call(job, downloadHandle)
	return
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
	return
}

func extractAudio(job *job.Job) {
	logger.Info.Printf("提取音频. [URL: %s] [JobID: %s] [Status: %d]\n", job.URL, job.JobID, job.Status)
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
		PackagePath: filepath.Join(util.WorkSpace, "audio_analysis"),
		FileName:    "api",
		MethodName:  "extract_audio",
		Args: []string{
			SetArg(filepath.Join(r.Location, r.VideoPath)),
		},
	}

	// 提取音频
	go python.Call(job, extractHandle)
	return
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
	job.SetStatus(util.JobExtractAudioDone)
	go JobSchedule(job)
	return
}

func waitDownload(job *job.Job) {
	logger.Info.Printf("资源已经存在, 等待下载完成. [URL: %s] [JobID: %s] [Status: %d]\n", job.URL, job.JobID, job.Status)
	// 获取资源信息
	r, err := resource.GetByKey(job.URL)
	if err != nil {
		logger.Error.Println("获取资源出错")
		job.SetStatus(util.JobErrFailedToFindResource)
		go JobSchedule(job)
		return
	}

	for {
		if r.Status == util.ResourceCompleted {
			logger.Info.Println("资源下载完成")
			job.SetStatus(util.JobExtractAudioDone)
			go JobSchedule(job)
			return
		} else if r.Status == util.ResourceErrDownloadFailed {
			logger.Error.Println("资源下载失败")
			job.SetStatus(util.JobErrDownloadFailed)
			go JobSchedule(job)
			return
		} else if r.Status == util.ResourceErrDownloadFailed {
			logger.Error.Println("音频提取失败")
			job.SetStatus(util.JobErrExtractFailed)
			go JobSchedule(job)
			return
		} else {
			logger.Warning.Println("资源下载中, 等待完成")
			time.Sleep(time.Second * 5)
			r.Refresh()
		}
	}
}

func extractAbstract(job *job.Job) {
	logger.Info.Printf("音频提取成功, 提取文本摘要和视频摘要. [URL: %s] [JobID: %s] [Status: %d]\n", job.URL, job.JobID, job.Status)
	// 判断资源是否已经存在
	if abstext.HaveAbsTextExisted(job.URL, job.KeyWords) {
		if job.Status&util.JobVideoAbstractExtractionDone != 0 {
			job.SetStatus(util.JobCompleted)
			go JobSchedule(job)
		} else {
			job.SetStatus(job.Status | util.JobTextAbstractExtractionDone)
		}
	} else {
		// 不存在就进行文本分析, 否则忽略
		go textAnalysis(job)
	}

	// 判断资源是否已经存在
	if absvideo.HaveAbsVideoExisted(job.URL) {
		if job.Status&util.JobTextAbstractExtractionDone != 0 {
			job.SetStatus(util.JobCompleted)
			go JobSchedule(job)
		} else {
			job.SetStatus(job.Status | util.JobVideoAbstractExtractionDone)
		}
	} else {
		// 不存在就进行视频分析, 否则忽略
		go videoAnalysis(job)
	}
}
