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

// PostJob is used to process "/job" post requests, job will be created
func PostJob(c *gin.Context) {
	// 获取数据
	json, err := util.GetJSON(c)
	if err != nil {
		logger.Warning.Println("POST数据解析失败")
		c.String(http.StatusBadRequest, err.Error())
		return
	}

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
		c.JSON(http.StatusBadRequest, gin.H{"warning": "Exists"})
		go server.JobSchedule(&oldJob)
		return
	}

	// 插入数据库
	err = mongodb.InsertOne(newJob)
	if err != nil {
		logger.Error.Printf("插入数据库失败. 原始数据: %+v, err:%+v\n", newJob, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "errorHappended"})
		return
	}

	// 开始下载
	logger.Info.Printf("接收到请求, 开始下载. [URL: %s] [JobID: %s]\n", newJob.URL, newJob.JobID)
	go server.JobSchedule(&newJob)

	// 返回 JobID
	c.JSON(http.StatusOK, gin.H{"job_id": newJob.JobID})
}
