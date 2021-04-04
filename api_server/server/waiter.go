package server

// func waitDownload(task *task.Task, r *resource.Resource) {
// 	logger.Info.Printf("[等待资源] 资源已经存在, 等待下载完成: %+v.\n", task)
// 	// 获取资源信息
// 	for {
// 		switch r.Status {
// 		case util.ResourceCompleted:
// 			logger.Info.Println("[等待资源] 资源下载完成, 音频提取成功")
// 			task.Status = util.JobToTextExtract
// 			redis.UpdataOne(task)
// 			return
// 		case util.ResourceErrDownloadFailed:
// 			logger.Error.Println("[等待资源] 资源下载失败")
// 			task.Status = util.JobErrDownloadFailed
// 			redis.UpdataOne(task)
// 			return
// 		case util.ResourceErrAudioExtractFailed:
// 			logger.Error.Println("[等待资源] 音频提取失败")
// 			task.Status = util.JobErrAudioExtractFailed
// 			redis.UpdataOne(task)
// 			return
// 		default:
// 			logger.Warning.Println("[等待资源] 资源下载中, 等待完成")
// 			time.Sleep(time.Second * 20)
// 		}

// 		// 获取最新状态
// 		redis.FindOne(r)
// 	}
// }
