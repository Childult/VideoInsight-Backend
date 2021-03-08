package router

import (
	"net/http"
	"swc/mongodb"
	"swc/mongodb/job"
	"swc/util"

	"github.com/gin-gonic/gin"
)

// GetJob is used to process "/media" post requests, deviceid will be return
func GetJob(c *gin.Context) {
	// 获取数据
	json, err := util.GetJSON(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	// 构建任务
	job := job.Job{
		JobID: json.GetID(),
	}

	// 查找数据
	data, err := mongodb.FindOne(job)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	status := data["status"]
	c.JSON(http.StatusOK, gin.H{"status": status})

}
