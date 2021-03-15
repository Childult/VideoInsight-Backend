package router

import (
	"io/ioutil"
	"net/http"
	"os"
	"swc/mongodb"
	"swc/mongodb/abstext"
	"swc/mongodb/absvideo"
	"swc/mongodb/job"
	"swc/mongodb/resource"
	"swc/util"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
	if job.Status == util.JobCompleted {
		at := abstext.AbsText{Hash: job.AbsText}
		text, err := mongodb.FindOne(at)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Not Found Text Abstract"})
			return
		}

		r, _ := resource.GetByKey(job.URL)
		av := absvideo.AbsVideo{URL: job.URL}
		video, err := mongodb.FindOne(av)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Not Found Video Abstract"})
			return
		}
		prefix := r.Location
		pictures := video["abstract"]
		pics := make(map[string]string)

		for _, x := range pictures.(bson.A) {
			file, _ := os.Open(prefix + x.(string))
			content, _ := ioutil.ReadAll(file)
			pics[x.(string)] = string(content)
		}

		c.JSON(http.StatusOK, gin.H{
			"text":     text,
			"pictures": pics,
		})
	}
	c.JSON(http.StatusOK, gin.H{"status": job.Status})

}
