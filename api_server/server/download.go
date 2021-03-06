package server

import (
	"fmt"
	"strconv"
	"swc/mongodb"
	"time"
)

const (
	completed = iota
	exist
	exception
)

func downloadMedia(job *mongodb.Job) int {
	// 构建资源
	source := mongodb.Source{
		URL:      job.URL,
		Status:   mongodb.Downloading,
		Location: "/home/download/" + strconv.FormatInt(time.Now().Unix(), 10) + "/",
	}

	// 检查资源是否存在
	exists := mongodb.HaveExisted(source)
	if exists {
		return exist
	}

	// 构建视频下载对象
	videoGetterPath := "/home/backend/SWC-Backend/video_getter/"
	fileName := "main"
	methodName := "download_video"
	args := []PyArgs{
		ArgsTemp(source.URL),
		ArgsTemp(source.Location),
	}

	python := PyWorker{
		PackagePath: videoGetterPath,
		FileName:    fileName,
		MethodName:  methodName,
		Args:        args,
	}

	// 调用python
	result := python.Call()
	if len(result) != 1 {
		return exception
	}

	// 下载成功
	source.VideoPath = result[0]

	// 构建音频提取对象
	python.PackagePath = "/home/backend/SWC-Backend/video_analysis/"
	python.FileName = "extract_audio"
	python.MethodName = "extract_audio"
	python.Args = []PyArgs{
		ArgsTemp(source.Location + source.VideoPath),
	}

	// 调用python
	result = python.Call()
	if len(result) != 1 {
		return exception
	}

	// 下载成功
	source.AudioPath = result[0]

	// 写入数据库
	source.Status = mongodb.Completed
	mongodb.Update(source)
	fmt.Println(source)

	job.Status = mongodb.Processing
	mongodb.Update(job)

	return completed
}
