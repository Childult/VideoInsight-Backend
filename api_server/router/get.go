package router

import (
	"net/http"
	"swc/mongodb"
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
	job := mongodb.Job{
		DeviceID: json.DeviceID,
		URL:      json.URL,
		KeyWords: json.KeyWords,
		JobID:    json.GetID(),
		Status:   mongodb.Downloading,
	}

	// 查找数据
	data, err := mongodb.FindOne(job)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	status := data["status"]
	if status != mongodb.Completed {
		// 返回
		c.JSON(http.StatusOK, gin.H{"status": status})
	} else {

	}

}
