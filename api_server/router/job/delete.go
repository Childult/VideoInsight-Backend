package job_router

import (
	"net/http"
	"swc/data/job"
	"swc/dbs/mongodb"
	"swc/dbs/redis"
	"swc/logger"

	"github.com/gin-gonic/gin"
)

// DeleteJob handle delete "/job"
var DeleteJob = func(c *gin.Context) {
	logger.Debug.Println("[DELETE] 开始")
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

	newJob := job.NewJob(jpm.DeviceID, jpm.URL, jpm.KeyWords)

	// 删除数据库
	if redis.Exists(newJob) || mongodb.Exists(newJob) {
		redis.DeleteOne(newJob)
		mongodb.DeleteOne(newJob)
		rt = ReturnType{
			Status:  0,
			Message: "删除成功",
			Result:  ""}
	} else {
		rt = ReturnType{
			Status:  -2,
			Message: "资源不存在",
			Result:  ""}
	}

	// 返回
	logger.Debug.Println("[DELETE] 结束")
	c.JSON(http.StatusBadRequest, rt)
}
