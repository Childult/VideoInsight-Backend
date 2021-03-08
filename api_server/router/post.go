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
	job := job.Job{
		DeviceID: json.DeviceID,
		URL:      json.URL,
		KeyWords: json.KeyWords,
		JobID:    json.GetID(),
		Status:   util.JobStart,
	}

	// 插入数据库
	err = mongodb.InsertOne(job)
	if err != nil {
		logger.Error.Printf("插入数据库失败. 原始数据: %+v\n", job)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// 开始下载
	logger.Info.Printf("接收到请求, 开始下载. [URL: %s] [JobID: %s]\n", job.URL, job.JobID)
	go server.JobSchedule(&job)

	// 返回 JobID
	c.JSON(http.StatusOK, gin.H{"job_id": job.JobID})
}
