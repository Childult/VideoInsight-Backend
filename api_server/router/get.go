package router

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"swc/logger"
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
	logger.Info.Println("[GET] 开始")
	var rt ReturnType
	// 获取数据
	jobID := c.Param("job_id")

	// 查找数据
	job, err := job.GetByKey(jobID)
	if err != nil {
		// 获取资源出错
		logger.Error.Println("[GET] 获取资源出错")
		rt = ReturnType{
			Status:  -3,
			Message: fmt.Sprintf("非找到`job_id=%s`的任务", jobID),
			Result:  ""}
		c.JSON(http.StatusBadRequest, rt)
		return
	}

	// 如果任务已经完成
	if job.Status == util.JobCompleted {
		at := abstext.AbsText{Hash: job.AbsText}
		text, _ := mongodb.FindOne(at)

		r, _ := resource.GetByKey(job.URL)
		av := absvideo.AbsVideo{URL: job.URL}
		video, _ := mongodb.FindOne(av)

		prefix := r.Location
		pictures := video["abstract"]
		pics := make(map[string]string)

		for _, x := range pictures.(bson.A) {
			file, _ := os.Open(prefix + x.(string))
			content, _ := ioutil.ReadAll(file)
			pics[x.(string)] = string(content)
		}

		rt = ReturnType{
			Status:  int(job.Status),
			Message: util.GetJobStatus(job.Status),
			Result: gin.H{
				"text":     text,
				"pictures": pics,
			}}
	} else {
		rt = ReturnType{
			Status:  int(job.Status),
			Message: util.GetJobStatus(job.Status),
			Result:  ""}
	}
	logger.Info.Printf("[GET] 返回状态%+v.\n", rt)
	c.JSON(http.StatusOK, rt)
}
