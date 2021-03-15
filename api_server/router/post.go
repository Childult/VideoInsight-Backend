package router

import (
	"net/http"
	"swc/logger"
	"swc/mongodb"
	"swc/mongodb/job"
	"swc/server"
	"swc/util"

	"github.com/gin-gonic/gin"
)

type ReturnType struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

// PostJob is used to process "/job" post requests, job will be created
func PostJob(c *gin.Context) {
	logger.Info.Println("[POST] 开始")
	var rt ReturnType

	// 获取数据
	json, err := util.GetJSON(c)
	if err != nil {
		logger.Error.Printf("[POST] 数据解析失败: %+v.\n", err)
		rt = ReturnType{
			Status:  -1,
			Message: err.Error(),
			Result:  ""}
		c.JSON(http.StatusBadRequest, rt)
		return
	}
	logger.Info.Printf("[POST] 收到: %+v.\n", json)

	// 构建任务
	newJob := job.Job{
		DeviceID: json.DeviceID,
		URL:      json.URL,
		KeyWords: json.KeyWords,
		JobID:    json.GetID(),
		Status:   util.JobStart,
	}

	// 查找数据
	oldJob, err := job.GetByKey(newJob.JobID)
	if err == nil {
		// 找到过去的数据
		logger.Warning.Printf("[POST] 该任务已存在: %+v.\n", oldJob)
		rt = ReturnType{
			Status:  0,
			Message: "该任务已存在",
			Result:  gin.H{"job_id": oldJob.JobID}}
		c.JSON(http.StatusOK, rt)
		go server.JobSchedule(&oldJob)
		return
	}

	// 插入数据库
	err = mongodb.InsertOne(newJob)
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
	go server.JobSchedule(&newJob)

	// 返回 JobID
	rt = ReturnType{
		Status:  0,
		Message: "任务创建成功",
		Result:  gin.H{"job_id": oldJob.JobID}}

	c.JSON(http.StatusOK, rt)
}
