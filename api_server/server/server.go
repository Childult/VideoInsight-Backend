package server

import (
	"fmt"
	"swc/mongodb"
	"time"
)

// StartTask .
func StartTask(job mongodb.Job) {
	go downloadMedia(&job)
}

func waiter(job *mongodb.Job) {
	// 构建资源对象
	source := mongodb.Source{
		URL: job.URL,
	}
	for {
		// 获取资源状态
		data, err := mongodb.FindOne(source)
		if err != nil {
			// 获取资源出错
			job.Status = mongodb.ErrorHappended
			mongodb.Update(job)
			break
		}
		// 资源下载完成
		if data["status"] == mongodb.Completed {
			source.Status = data["status"].(string)
			source.Location = data["location"].(string)
			source.VideoPath = data["videopath"].(string)
			source.AudioPath = data["audiopath"].(string)

			fmt.Println("---------  进行文本分析 ------------------")
			// 进行文本分析
			textAnalysis(job, source)
			return
			// 进行视频分析
			// break
		} else if data["status"] == mongodb.Downloading {
			fmt.Println("---------  等待 ------------------")
			time.Sleep(time.Second * 5)
		} else {
			fmt.Println("---------  文本异常 ------------------")
			fmt.Println(data["status"])
			break
		}
	}
}
