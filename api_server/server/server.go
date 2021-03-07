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
	// Start 创建资源, 写入数据库
	Start = iota
	// DownloadMedia 下载资源
	DownloadMedia
	// ExtractAudio 提取音频
	ExtractAudio
	// Existed 资源已存在
	Existed
	// TaskErr 错误发生
	TaskErr
)

// PythonHandlerFunc python 回调函数
type PythonHandlerFunc func(job job.Job, result string)

// Schedule .
func Schedule(status int, job job.Job) {
	switch status {
	case Start:
		go creatSource(job)

	case DownloadMedia:
		go mediaDownload(job)

	case ExtractAudio:

	case Existed:

	case TaskErr:

	}
}

// TaskSchedule .
func TaskSchedule(status int, job job.Job) {
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
