package test_router

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"swc/data/abstext"
	"swc/data/absvideo"
	"swc/data/job"
	"swc/data/resource"
	"swc/logger"
	"swc/util"

	"github.com/gin-gonic/gin"
)

// GetJob is used to process "/job" post requests, deviceid will be return
func GetJob(c *gin.Context) {
	logger.Debug.Println("[GET] 开始")
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

	// 查找数据
	newJob := job.NewJob(jpm.DeviceID, jpm.URL, jpm.KeyWords)
	err = newJob.Load()
	if err == nil {
		// 任务已完成
		r := resource.Resource{URL: newJob.URL}
		r.Load()

		at := abstext.NewAbsText(newJob.URL, newJob.KeyWords)
		at.Load()

		av := absvideo.AbsVideo{URL: newJob.URL}
		av.Load()

		prefix := r.Location
		pictures := make(map[string]string)

		for _, filename := range av.Abstract {
			file, _ := os.Open(prefix + filename)
			content, _ := ioutil.ReadAll(file)
			pictures[filename] = string(content)
		}

		rt = ReturnType{
			Status:  int(newJob.Status),
			Message: util.GetJobStatus(newJob.Status),
			Result: gin.H{
				"text":     at.Abstract,
				"pictures": pictures,
			}}
	}

	err = newJob.Retrieve()
	if err != nil {
		// 获取资源出错
		logger.Error.Println("[GET] 获取资源出错")
		rt = ReturnType{
			Status:  -3,
			Message: fmt.Sprintf("未找到`job_id=%s`的任务", newJob.JobID),
			Result:  ""}
		c.JSON(http.StatusBadRequest, rt)
		return
	}

	// 如果任务已经完成
	rt = ReturnType{
		Status:  int(newJob.Status),
		Message: util.GetJobStatus(newJob.Status),
		Result:  ""}

	logger.Debug.Printf("[GET] 返回状态%+v.\n", rt)
	c.JSON(http.StatusOK, rt)
}
