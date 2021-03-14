package server

import (
	"swc/logger"
	"swc/mongodb/job"
	"swc/util"
)

// JobSchedule .
func JobSchedule(job *job.Job) {
	status := job.Status

	if status > util.JobCompleted {
		logger.Error.Println("任务失败, 错误代码:", status)
		return
	}

	if status == util.JobCompleted {
		logger.Info.Printf("任务完成. [URL: %s] [JobID: %s] [Status: %d]\n", job.URL, job.JobID, job.Status) // 任务完成

	} else if status == util.JobStart {
		go creatResource(job) // 创建资源

	} else if status == util.JobDownloadMedia {
		go mediaDownload(job) // 下载视频

	} else if status == util.JobExisted {
		go waitDownload(job) // 资源已经存在, 等待下载完成

	} else if status == util.JobExtractAudio {
		go extractAudio(job) // 提取音频

	} else if status == util.JobExtractAudioDone {
		go extractAbstract(job) // 音频提取成功, 提取文本摘要和视频摘要

	} else if status&util.JobVideoAbstractExtractionDone != 0 {
		go textAnalysis(job) // 视频摘要完成, 但文本摘要未完成

	} else if status&util.JobTextAbstractExtractionDone != 0 {
		go videoAnalysis(job) // 文本摘要完成, 但视频摘要未完成

	}
}
