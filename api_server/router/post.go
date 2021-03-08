package router

import (
	"fmt"
	"net/http"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// 开始下载
	fmt.Println("========================= 接收到请求, 开始下载 ===================================")
	go server.JobSchedule(&job)

	// 返回 JobID
	c.JSON(http.StatusOK, gin.H{"jobid": job.JobID})
}
