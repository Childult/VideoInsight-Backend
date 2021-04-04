package job_router

import (
	"net/http"
	"swc/data/job"
	"swc/dbs/mongodb"
	"swc/logger"
	"swc/server/task_builder"
	"sync"

	"github.com/gin-gonic/gin"
)

type ReturnType struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

// JobPostMessage 用户利用 POST 提交的数据, 用于为任务创建唯一的 ID
type JobPostMessage struct {
	DeviceID string   `json:"device_id"`
	URL      string   `json:"url"`
	KeyWords []string `json:"keywords,omitempty"`
}

var jobMu sync.Mutex

// PostJob 接收并处理发往 /job 的 POST 请求
var PostJob = func(c *gin.Context) {
	logger.Info.Println("[POST Job] 收到 POST 请求")
	var rt ReturnType
	var err error

	// 解析 json 字段
	var jpm JobPostMessage
	err = c.ShouldBindJSON(&jpm)
	if err != nil {
		logger.Error.Printf("[POST Job] 数据解析失败: %+v.\n", err)
		rt = ReturnType{
			Status:  jsonErr,
			Message: err.Error(),
			Result:  ""}
		c.JSON(http.StatusBadRequest, rt)
		return
	}

	logger.Info.Printf("[POST Job] 开始创建任务: %+v.\n", jpm)

	// 构建任务
	newJob := job.NewJob(jpm.DeviceID, jpm.URL, jpm.KeyWords)
	jobMu.Lock()
	if mongodb.Exists(newJob) {
		jobMu.Unlock()
		// 返回 JobID
		rt = ReturnType{
			Status:  0,
			Message: "任务已存在",
			Result:  gin.H{"job_id": newJob.JobID}}
	} else {
		mongodb.InsertOne(newJob)
		jobMu.Unlock()

		// 开启任务
		go task_builder.AddTask(newJob.URL, newJob.KeyWords)

		// 返回 JobID
		rt = ReturnType{
			Status:  0,
			Message: "任务已接收",
			Result:  gin.H{"job_id": newJob.JobID}}
	}
	c.JSON(http.StatusOK, rt)
}
