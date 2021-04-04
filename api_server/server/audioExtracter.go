package server

import (
	"path/filepath"
	"strings"
	"swc/data/resource"
	"swc/data/task"
	"swc/dbs/redis"
	"swc/logger"
	"swc/util"
)

func extractAudio(job *task.Task, r *resource.Resource) {
	logger.Info.Printf("[提取音频] 开始: %+v.\n", job)
	// 开始提取音频
	r.Status = util.ResourceAudioExtracting
	redis.UpdataOne(r)
	// 构建音频提取对象
	python := PyWorker{
		PackagePath: filepath.Join(util.WorkSpace, "audio_analysis"), // 包名
		FileName:    "api",                                           // 文件名
		MethodName:  "extract_audio",                                 // 调用函数
		Args: []string{ // 实参
			SetArg(filepath.Join(r.Location, r.VideoPath)), // 传入视频
		},
	}

	// 提取音频
	result := python.Call()
	// 是否成功提取音频
	if len(result) == 0 {
		logger.Error.Println("[提取音频回调] 提取失败.")
		r.Status = util.ResourceErrAudioExtractFailed
		job.Status = util.JobErrAudioExtractFailed
	} else {
		// 音频提取成功, 更新状态
		pythonReturn := strings.Join(result, "")
		logger.Debug.Printf("[提取音频回调] 音频提取成功: %+v.\n", pythonReturn)
		r.AudioPath = pythonReturn

		r.Status = util.ResourceCompleted
		job.Status = util.JobAudioExtractDone
	}
	redis.UpdataOne(r)
	redis.UpdataOne(job)
}
