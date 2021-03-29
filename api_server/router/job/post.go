package job_router

import (
	"net/http"
	"swc/data/job"
	"swc/logger"
	"swc/server"

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

// PostJob is used to process "/job" post requests, job will be created
func PostJob(c *gin.Context) {
	logger.Info.Println("[POST] 开始")
	var rt ReturnType
	var err error

	// 获取数据
	var jpm JobPostMessage
	err = c.ShouldBindJSON(&jpm)
	if err != nil {
		logger.Error.Printf("[POST] 数据解析失败: %+v.\n", err)
		rt = ReturnType{
			Status:  -1,
			Message: err.Error(),
			Result:  ""}
		c.JSON(http.StatusBadRequest, rt)
		return
	}

	logger.Debug.Printf("[POST] 收到: %+v.\n", jpm)

	// 构建任务
	newJob := job.NewJob(jpm.DeviceID, jpm.URL, jpm.KeyWords)

	// 查找数据
	oldJob := job.Job{}
	oldJob.JobID = newJob.JobID
	err = oldJob.Load()
	if err == nil {
		// 任务已完成
		logger.Warning.Printf("[POST] 该任务已完成: %+v.\n", oldJob)
		rt = ReturnType{
			Status:  0,
			Message: "该任务已完成",
			Result:  gin.H{"job_id": oldJob.JobID}}
		c.JSON(http.StatusOK, rt)
		return
	}

	err = oldJob.Retrieve()
	if err == nil {
		// 任务未完成, 但已经提交
		logger.Warning.Printf("[POST] 该任务已存在: %+v.\n", oldJob)
		rt = ReturnType{
			Status:  0,
			Message: "该任务已存在",
			Result:  gin.H{"job_id": oldJob.JobID}}
		c.JSON(http.StatusOK, rt)
		return
	}

	// 插入数据库
	err = newJob.Save()
	if err != nil {
		logger.Error.Printf("[POST] 插入数据库失败. 原始数据: %+v, err:%+v\n", newJob, err.Error())
		rt = ReturnType{
			Status:  -2,
			Message: "插入数据库失败",
			Result:  ""}
		c.JSON(http.StatusBadRequest, rt)
		return
	}

	// 开始下载
	logger.Info.Printf("[POST] 接收到请求, 开始执行任务: %+v.\n", newJob)
	go server.JobSchedule(newJob)

	// 返回 JobID
	rt = ReturnType{
		Status:  0,
		Message: "任务创建成功",
		Result:  gin.H{"job_id": newJob.JobID}}

	c.JSON(http.StatusOK, rt)
}
