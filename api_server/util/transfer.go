package util

func GetTaskStatus(status int32) string {
	switch status {
	case TaskStart:
		return "任务已接收, 即将开始"

	case TaskCreated:
		return "任务进行中"

	case TaskDownloadDone:
		return "资源下载完成, 即将开始提取摘要"

	case TaskCompleted:
		return "任务已完成"

	default:
		return "出现错误, 请删除任务后重试"
	}
}
