package server

import (
	"fmt"
	"swc/mongodb"
	"swc/mongodb/job"
	"swc/mongodb/source"
	"swc/util"
	"time"
)

const (
	Start = iota
	DownloadMedia
	TaskErr
)

// TaskSchedule .
func TaskSchedule(status int, job job.Job) {
	// switch status {
	// case Start:
	// 	creatSource(job)
	// }
	go downloadMedia(&job)
}

func waiter(job *job.Job) {
	// 构建资源对象
	source := source.Source{
		URL: job.URL,
	}
	for {
		// 获取资源状态
		err := source.Refresh()
		if err != nil {
			// 获取资源出错
			job.Status = util.ErrorHappended
			mongodb.Update(job)
			break
		}
		// 资源下载完成
		if source.Status == util.Completed {

			fmt.Println("---------  进行文本分析 ------------------")
			// 进行文本分析
			textAnalysis(job, source)
			// 进行视频分析
			break
		} else if source.Status == util.Downloading {
			fmt.Println("---------  等待 ------------------")
			time.Sleep(time.Second * 5)
		} else {
			fmt.Println("---------  文本异常 ------------------")
			fmt.Println(source.Status)
			break
		}
	}
}
