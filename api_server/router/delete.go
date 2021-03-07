package router

import (
	"net/http"
	"swc/mongodb"
	"swc/util"

	"github.com/gin-gonic/gin"
)

// DeleteJob handle delete "/job"
func DeleteJob(c *gin.Context) {
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
	// 删除数据库
	err = mongodb.DeleteOne(job)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 返回
	c.String(http.StatusOK, "")

}
