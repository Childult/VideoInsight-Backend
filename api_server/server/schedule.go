package server

import (
	"swc/data/job"
	"swc/data/resource"
	"swc/data/task"
	"swc/dbs/mongodb"
	"swc/dbs/redis"
	"swc/logger"
	"swc/util"
	"sync"
)

func StartJob(job *job.Job) {
	// 目前不考虑有没有关键字, 直接创建
	t := task.NewTask(job.URL, job.KeyWords)
	StartTask(t)
}

var taskMu sync.Mutex

func StartTask(task *task.Task) {
	// 加锁, 保证一个任务不会被创建两次. 一个任务对应一个URL+关键字
	taskMu.Lock()
	defer taskMu.Unlock()

	// 先查看该任务是否已经完成
	if redis.Exists(task) || mongodb.Exists(task) {
		// 重复提交任务, 该任务存在
		logger.Warning.Printf("[Start task] 重复提交任务, 该任务存在: %+v.\n", task)
	} else {
		// 都不存在, 创建一个新的任务
		err := redis.InsertOne(task)
		if err != nil {
			logger.Error.Printf("[Start task] 插入数据库失败. 原始数据: <%+v>, error: <%s>\n", task, err)
			return
		}

		// 任务创建成功, 开始执行
		logger.Info.Printf("[Start task] 任务创建成功, 开始执行: %+v.\n", task)
		go taskSchedule(task)
	}
}

// taskSchedule .
func taskSchedule(task *task.Task) {
	var r *resource.Resource
	wg := new(sync.WaitGroup)

	// 不断循环, 直到任务完成或发生错误
	for task.Status < util.JobCompleted {
		switch task.Status {
		// 创建资源
		case util.JobStart:
			r = creatResource(task)

		// 下载资源
		case util.JobToDownloadMedia:
			mediaDownload(task, r)

		// 提取音频
		case util.JobToExtractAudio:
			extractAudio(task, r)

		// 只进行文本分析, 视频分析会在另一个协程里完成
		case util.JobToTextExtract:
			redis.FindOne(r)
			wg.Add(1)
			go extractTextAbstract(task, r) // 音频提取成功, 提取文本摘要
			wg.Wait()

		// 音频提取成功, 开始提取文本摘要和视频摘要
		case util.JobAudioExtractDone:
			wg.Add(1)
			go extractTextAbstract(task, r) // 音频提取成功, 提取文本摘要
			wg.Add(1)
			go extractVideoAbstract(task, r) // 音频提取成功, 提取视频摘要
			wg.Wait()
		}
	}
	if task.Status == util.JobCompleted {
		logger.Info.Printf("[任务调度] 任务完成: %+v.\n", task)
	} else if task.Status > util.JobCompleted {
		logger.Error.Printf("[任务调度] 任务失败, 错误代码: %d, 原因: %s.\n", task.Status, util.GetJobStatus(task.Status))
	}
}

var absMu sync.Mutex

func setAbstractFlag(t *task.Task, status int32) {
	absMu.Lock()
	defer absMu.Unlock()
	if t.Status > util.JobCompleted {
		return
	} else if t.Status == util.JobTextAbstractExtractionDone && status == util.JobVideoAbstractExtractionDone ||
		t.Status == util.JobVideoAbstractExtractionDone && status == util.JobTextAbstractExtractionDone {
		t.Status = util.JobCompleted
		redis.UpdataOne(t)
		mongodb.InsertOne(t)
	} else {
		t.Status = status
		redis.UpdataOne(t)
	}
}
