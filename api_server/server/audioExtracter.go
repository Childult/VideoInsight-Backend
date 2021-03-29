package server

import (
	"path/filepath"
	"strings"
	"swc/data/job"
	"swc/data/resource"
	"swc/logger"
	"swc/util"
)

func extractAudio(job *job.Job) {
	logger.Info.Printf("[提取音频] 开始: %+v.\n", job)
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
	r := resource.Resource{URL: job.URL}
	err := r.Retrieve()
	if err != nil {
		logger.Error.Printf("[提取音频回调] 获取资源出错: %+v.\n", err)
		job.Status = util.JobErrFailedToFindResource
		job.Save()
		go JobSchedule(job)
		return
	}

	// 是否成功提取音频
	if len(result) == 0 {
		logger.Error.Println("[提取音频回调] 提取失败.")
		r.Status = util.ResourceErrExtractFailed
		r.Save()
		job.Status = util.JobErrExtractFailed
		job.Save()
		go JobSchedule(job)
		return
	}
	pythonReturn := strings.Join(result, "")

	// 音频提取成功, 更新状态
	logger.Debug.Printf("[提取音频回调] 音频提取成功: %+v.\n", pythonReturn)
	r.AudioPath = pythonReturn
	r.Status = util.ResourceCompleted
	r.Save()
	job.Status = util.JobExtractAudioDone
	job.Save()
	go JobSchedule(job)
}
