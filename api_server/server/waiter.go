package server

import (
	"swc/data/job"
	"swc/data/resource"
	"swc/logger"

	"swc/util"
	"time"
)

func waitDownload(job *job.Job) {
	logger.Info.Printf("[候车间] 资源已经存在, 等待下载完成: %+v.\n", job)
	// 获取资源信息
	r := resource.Resource{URL: job.URL}
	err := r.Retrieve()
	if err != nil {
		logger.Error.Printf("[候车间] 获取资源出错: %+v.\n", err)
		job.Status = util.JobErrFailedToFindResource
		job.Save()
		go JobSchedule(job)
		return
	}

	for {
		logger.Debug.Printf("[候车间] 检查资源状态: %+v.\n", r)
		if r.Status == util.ResourceCompleted {
			logger.Debug.Println("[候车间] 资源下载完成")
			job.Status = util.JobExtractAudioDone
			job.Save()
			go JobSchedule(job)
			return
		} else if r.Status == util.ResourceErrDownloadFailed {
			logger.Error.Println("[候车间] 资源下载失败")
			job.Status = util.JobErrDownloadFailed
			job.Save()
			go JobSchedule(job)
			return
		} else if r.Status == util.ResourceErrExtractFailed {
			logger.Error.Println("[候车间] 音频提取失败")
			job.Status = util.JobErrExtractFailed
			job.Save()
			go JobSchedule(job)
			return
		} else {
			logger.Warning.Println("[候车间] 资源下载中, 等待完成")
			time.Sleep(time.Second * 20)
			err = r.Retrieve()
			if err != nil {
				logger.Error.Printf("[候车间] 获取资源出错: %+v.\n", err)
				job.Status = util.JobErrFailedToFindResource
				job.Save()
				go JobSchedule(job)
				return
			}
		}
	}
}
