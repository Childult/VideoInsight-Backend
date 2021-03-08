package server

import (
	"swc/logger"
	"swc/mongodb/job"
	"swc/util"
)

// JobSchedule .
func JobSchedule(job *job.Job) {
	switch job.Status {
	case util.JobStart: // 创建资源
		go creatResource(job)

	case util.JobDownloadMedia: // 下载视频
		go mediaDownload(job)

	case util.JobExtractAudio: // 提取音频
		go extractAudio(job)

	case util.JobExtractDone: // 音频提取成功, 提取文本摘要和视频摘要

		go extractAbstract(job)
	case util.JobAbstractextraction:

	case util.JobCompleted: // 任务完成
		logger.Info.Printf("任务完成. [URL: %s] [JobID: %s] [Status: %d]\n", job.URL, job.JobID, job.Status)

	case util.JobExisted: // 资源已经存在, 等待下载完成
		go waitDownload(job)

	case util.JobErrFailedToFindResource:
	}
}
