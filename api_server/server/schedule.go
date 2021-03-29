package server

import (
	"swc/data/job"
	"swc/logger"
	"swc/util"
)

// JobSchedule .
func JobSchedule(job *job.Job) {
	status := job.Status

	if status > util.JobCompleted {
		logger.Error.Printf("[任务调度] 任务失败, 错误代码: %d, 原因: %s.\n", status, util.GetJobStatus(status))
		return
	}

	if status&util.JobCompleted != 0 {
		logger.Info.Printf("[任务调度] 任务完成: %+v.\n", job)

	} else if status&util.JobStart != 0 {
		go creatResource(job) // 创建资源

	} else if status&util.JobDownloadMedia != 0 {
		go mediaDownload(job) // 下载视频

	} else if status&util.JobExisted != 0 {
		go waitDownload(job) // 资源已经存在, 等待下载完成

	} else if status&util.JobExtractAudio != 0 {
		go extractAudio(job) // 提取音频

	} else if status&util.JobExtractAudioDone != 0 {
		go extractTextAbstract(job)  // 音频提取成功, 提取文本摘要
		go extractVideoAbstract(job) // 音频提取成功, 提取视频摘要

	} else if status&util.JobTextAbstractExtractionDone != 0 {
		go extractVideoAbstract(job) // 文本摘要完成, 但视频摘要未完成

	} else if status&util.JobVideoAbstractExtractionDone != 0 {
		go extractTextAbstract(job) // 视频摘要完成, 但文本摘要未完成

	}
}
