package router

import (
	"net/http"
	"swc/mongodb/job"

	"github.com/gin-gonic/gin"
)

// GetJob is used to process "/job" post requests, deviceid will be return
func GetJob(c *gin.Context) {
	// 获取数据
	jobID := c.Param("job_id")

	// 查找数据
	job, err := job.GetByKey(jobID)
	if err != nil {
		// 获取资源出错
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not Found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": job.Status})

}
