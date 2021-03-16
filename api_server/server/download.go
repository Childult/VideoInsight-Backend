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
	logger.Info.Printf("[创建资源] 开始: %+v.\n", job)
	// 构建资源
	resource := resource.Resource{
		URL:      job.URL,                                                                      // 资源链接
		Status:   util.ResourceDownloading,                                                     // 状态: 下载中
		Location: filepath.Join(util.Location, strconv.FormatInt(time.Now().Unix(), 10)) + "/", // 保存地址, 以时间戳为文件夹
	}

	// 检查资源是否存在
	exists := mongodb.HaveExisted(resource)
	if exists {
		logger.Warning.Println("[创建资源] 资源已存在")
		job.SetStatus(util.JobExisted)
		go JobSchedule(job)
		return
	}

	// 首次写入数据库
	logger.Debug.Printf("[创建资源] 资源创建成功:%+v.\n", resource)
	mongodb.InsertOne(resource)
	job.SetStatus(util.JobDownloadMedia)
	go JobSchedule(job)
}

func mediaDownload(job *job.Job) {
	logger.Info.Printf("[下载视频] 开始: %+v.\n", job)
	// 获取资源信息
	resource, err := resource.GetByKey(job.URL)
	if err != nil {
		logger.Error.Printf("[下载视频] 获取资源出错: %+v.\n", err)
		job.SetStatus(util.JobErrFailedToFindResource)
		go JobSchedule(job)
		return
	}

	// 构建视频下载对象
	python := PyWorker{
		PackagePath: filepath.Join(util.WorkSpace, "video_getter"), // python 包地址
		FileName:    "api",                                         // 文件名
		MethodName:  "download_video",                              // 调用函数
		Args: []string{ // 实参
			SetArg(resource.URL),      // 资源链接
			SetArg(resource.Location), // 保存路径
		},
	}

	go python.Call(job, downloadHandle)
}

func downloadHandle(job *job.Job, result []string) {
	logger.Info.Printf("[下载视频回调] 开始: %+v.\n", job)
	// 获取资源信息
	r, err := resource.GetByKey(job.URL)
	if err != nil {
		logger.Error.Printf("[下载视频回调] 获取资源出错: %+v.\n", err)
		job.SetStatus(util.JobErrFailedToFindResource)
		go JobSchedule(job)
		return
	}

	// 是否下载成功
	if len(result) != 1 {
		logger.Error.Println("[下载视频回调] 下载失败.")
		r.SetStatus(util.ResourceErrDownloadFailed)
		job.SetStatus(util.JobErrDownloadFailed)
		go JobSchedule(job)
		return
	}

	// 下载成功, 更新状态
	logger.Debug.Printf("[下载视频回调] 视频下载成功: %+v.\n", result[0])
	r.VideoPath = result[0]
	r.SetStatus(util.ResourceExtracting)
	job.SetStatus(util.JobExtractAudio)
	go JobSchedule(job)
}

func extractAudio(job *job.Job) {
	logger.Info.Printf("[提取音频] 开始: %+v.\n", job)
	// 获取资源信息
	r, err := resource.GetByKey(job.URL)
	if err != nil {
		logger.Error.Printf("[下载视频回调] 获取资源出错: %+v.\n", err)
		job.SetStatus(util.JobErrFailedToFindResource)
		go JobSchedule(job)
		return
	}

	// 构建音频提取对象
	python := PyWorker{
		PackagePath: filepath.Join(util.WorkSpace, "audio_analysis"), // 包名
		FileName:    "api",                                           // 文件名
		MethodName:  "extract_audio",                                 // 调用函数
		Args: []string{ // 实参
			SetArg(filepath.Join(r.Location, r.VideoPath)), // 传入视频
		},
	}

	// 提取音频
	go python.Call(job, extractHandle)
}

func extractHandle(job *job.Job, result []string) {
	logger.Info.Printf("[提取音频回调] 开始: %+v.\n", job)
	// 获取资源信息
	r, err := resource.GetByKey(job.URL)
	if err != nil {
		logger.Error.Printf("[提取音频回调] 获取资源出错: %+v.\n", err)
		job.SetStatus(util.JobErrFailedToFindResource)
		go JobSchedule(job)
		return
	}

	// 是否成功提取音频
	if len(result) != 1 {
		logger.Error.Println("[提取音频回调] 提取失败.")
		r.SetStatus(util.ResourceErrExtractFailed)
		job.SetStatus(util.JobErrExtractFailed)
		go JobSchedule(job)
		return
	}

	// 音频提取成功, 更新状态
	logger.Debug.Printf("[提取音频回调] 音频提取成功: %+v.\n", result[0])
	r.AudioPath = result[0]
	r.SetStatus(util.ResourceCompleted)
	job.SetStatus(util.JobExtractAudioDone)
	go JobSchedule(job)
}

func waitDownload(job *job.Job) {
	logger.Info.Printf("[候车间] 资源已经存在, 等待下载完成: %+v.\n", job)
	// 获取资源信息
	r, err := resource.GetByKey(job.URL)
	if err != nil {
		logger.Error.Printf("[候车间] 获取资源出错: %+v.\n", err)
		job.SetStatus(util.JobErrFailedToFindResource)
		go JobSchedule(job)
		return
	}

	for {
		logger.Debug.Printf("[候车间] 检查资源状态: %+v.\n", r)
		if r.Status == util.ResourceCompleted {
			logger.Debug.Println("[候车间] 资源下载完成")
			job.SetStatus(util.JobExtractAudioDone)
			go JobSchedule(job)
			return
		} else if r.Status == util.ResourceErrDownloadFailed {
			logger.Error.Println("[候车间] 资源下载失败")
			job.SetStatus(util.JobErrDownloadFailed)
			go JobSchedule(job)
			return
		} else if r.Status == util.ResourceErrDownloadFailed {
			logger.Error.Println("[候车间] 音频提取失败")
			job.SetStatus(util.JobErrExtractFailed)
			go JobSchedule(job)
			return
		} else {
			logger.Warning.Println("[候车间] 资源下载中, 等待完成")
			time.Sleep(time.Second * 20)
			r.Refresh()
		}
	}
}

func extractTextAbstract(job *job.Job) {
	logger.Info.Printf("[提取文本摘要] 开始: %+v.\n", job)
	// 判断资源是否已经存在
	if abstext.HaveAbsTextExisted(job.URL, job.KeyWords) { // 文本摘要已经提取完成
		if job.Status&util.JobVideoAbstractExtractionDone != 0 { // 如果视频摘要已经提取完成
			job.SetStatus(util.JobCompleted) // 任务完成
			go JobSchedule(job)
		} else {
			job.SetStatus(util.JobTextAbstractExtractionDone)
			go JobSchedule(job)
		}
	} else {
		// 不存在就进行文本分析, 否则忽略
		go textAnalysis(job)
	}
}

func extractVideoAbstract(job *job.Job) {
	logger.Info.Printf("[提取视频摘要] 开始: %+v.\n", job)
	// 判断资源是否已经存在
	if absvideo.HaveAbsVideoExisted(job.URL) { // 视频摘要已经完成
		if job.Status&util.JobTextAbstractExtractionDone != 0 { // 如果文本摘要也已经完成
			job.SetStatus(util.JobCompleted) // 任务完成
			go JobSchedule(job)
		} else {
			job.SetStatus(util.JobVideoAbstractExtractionDone)
			go JobSchedule(job)
		}
	} else {
		// 不存在就进行视频分析, 否则忽略
		go videoAnalysis(job)
	}
}
