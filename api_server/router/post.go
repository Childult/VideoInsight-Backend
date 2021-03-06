package router

import (
	"net/http"
	"swc/mongodb"
	"swc/server"

	"github.com/gin-gonic/gin"
)

// PostJob is used to process "/job" post requests, job will be created
func PostJob(c *gin.Context) {
	// 获取数据
	json, err := getJSON(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	// 构建任务
	job := mongodb.Job{
		DeviceID: json.DeviceID,
		URL:      json.URL,
		KeyWords: json.KeyWords,
		JobID:    json.GetID(),
		Status:   mongodb.Downloading,
	}

	// 插入数据库
	err = mongodb.InsertOne(job)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// 开始下载
	go server.StartTask(job)

	// 返回 JobID
	c.JSON(http.StatusOK, gin.H{"jobid": job.JobID})
}
