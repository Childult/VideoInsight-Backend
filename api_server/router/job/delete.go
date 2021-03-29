package job_router

import (
	"net/http"
	"swc/data/job"
	"swc/logger"

	"github.com/gin-gonic/gin"
)

// DeleteJob handle delete "/job"
func DeleteJob(c *gin.Context) {
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
	err = newJob.Delete()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = newJob.Remove()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 返回
	logger.Debug.Println("[DELETE] 结束")
	c.String(http.StatusOK, "")

}
