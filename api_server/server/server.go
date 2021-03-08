package server

import (
	"fmt"
	"swc/mongodb/job"
	"swc/util"
)

// JobSchedule .
func JobSchedule(job *job.Job) {
	switch job.Status {
	case util.JobStart: // 创建资源
		fmt.Println("========================= 创建资源 ===================================")
		go creatResource(job)
	case util.JobDownloadMedia: // 下载视频
		fmt.Println("========================= 下载视频 ===================================")
		go mediaDownload(job)
	case util.JobExtractAudio: // 提取音频
		fmt.Println("========================= 提取音频 ===================================")
		go extractAudio(job)
	case util.JobExtractDone: // 音频提取成功, 提取文本摘要和视频摘要
		fmt.Println("========================= 音频提取成功, 提取文本摘要和视频摘要 ===================================")
		go extractAbstract(job)
	case util.JobAbstractextraction:

	case util.JobCompleted: // 任务完成
		fmt.Println("========================= 任务完成 ===================================")

	case util.JobExisted: // 资源已经存在, 等待下载完成
		fmt.Println("========================= 资源已经存在, 等待下载完成 ===================================")
		go waitDownload(job)
	case util.JobErrFailedToFindResource:
	}
}
