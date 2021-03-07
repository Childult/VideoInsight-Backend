package server

import (
	"fmt"
	"path/filepath"
	"strconv"
	"swc/mongodb"
	"swc/mongodb/job"
	"swc/mongodb/source"
	"swc/util"
	"time"
)

func creatSource(job job.Job) {
	// 构建资源
	source := source.Source{
		URL:      job.URL,
		Status:   util.Downloading,
		Location: filepath.Join(util.SavePath, strconv.FormatInt(time.Now().Unix(), 10)),
	}

	// 检查资源是否存在
	exists := mongodb.HaveExisted(source)
	if exists {
		go TaskSchedule(TaskErr, job)
		return
	}

	// 首次写入数据库
	mongodb.InsertOne(source)
	go TaskSchedule(DownloadMedia, job)
}

func downloadMedia(job *job.Job) {
	// 构建资源
	source := source.Source{
		URL:      job.URL,
		Status:   util.Downloading,
		Location: filepath.Join(util.SavePath, strconv.FormatInt(time.Now().Unix(), 10)),
	}

	// 检查资源是否存在
	exists := mongodb.HaveExisted(source)
	if exists {
		go waiter(job)
		return
	}
	// 首次写入数据库
	mongodb.InsertOne(source)

	// 构建视频下载对象
	videoGetterPath := filepath.Join(util.WorkSpace, "video_getter")
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

	// 开始下载
	fmt.Println("=====================================================")
	fmt.Println(python)
	result := python.Call()
	if len(result) != 1 {
		// 异常
		source.Status = util.ErrorHappended
		mongodb.Update(source)
		go waiter(job)
		return
	}

	// 下载成功, 更新数据库状态
	source.VideoPath = result[0]
	source.Status = util.Processing
	mongodb.Update(source)

	// 构建音频提取对象
	python.PackagePath = filepath.Join(util.WorkSpace, "video_analysis")
	python.FileName = "extract_audio"
	python.MethodName = "extract_audio"
	python.Args = []PyArgs{
		ArgsTemp(filepath.Join(source.Location, source.VideoPath)),
	}

	// 提取音频
	result = python.Call()
	if len(result) != 1 {
		// 异常
		source.Status = util.ErrorHappended
		mongodb.Update(source)
		go waiter(job)
		return
	}

	// 音频提取成功, 更新数据库
	source.AudioPath = result[0]
	source.Status = util.Completed
	mongodb.Update(source)

	// 更新任务状态
	job.Status = util.Processing
	mongodb.Update(job)
	go waiter(job)
}

func extractAudio() {

}
