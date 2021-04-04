package job_router

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"swc/data/abstext"
	"swc/data/absvideo"
	"swc/data/job"
	"swc/data/resource"
	"swc/data/task"
	"swc/dbs/mongodb"
	"swc/logger"
	"swc/util"

	"github.com/gin-gonic/gin"
)

const (
	jsonErr      = -iota - 1 // JSON 格式错误
	jobNotExists             // 任务不存在
)

// GetJob is used to process "/job" post requests, deviceid will be return
var GetJob = func(c *gin.Context) {
	logger.Debug.Println("[GET] 开始")
	var rt ReturnType
	var err error

	// 获取数据
	var jpm JobPostMessage
	err = c.ShouldBindJSON(&jpm)
	if err != nil {
		logger.Error.Printf("[POST] 数据解析失败: %+v.\n", err)
		rt = ReturnType{
			Status:  jsonErr,
			Message: err.Error(),
			Result:  ""}
		c.JSON(http.StatusBadRequest, rt)
		return
	}

	j := job.NewJob(jpm.DeviceID, jpm.URL, jpm.KeyWords)
	if !mongodb.Exists(j) {
		// 获取资源出错
		logger.Warning.Println("[GET] 查询了不存在的任务")
		rt = ReturnType{
			Status:  jobNotExists,
			Message: fmt.Sprintf("未找到`job_id=%s`的任务", j.JobID),
			Result:  ""}
		c.JSON(http.StatusBadRequest, rt)
		return
	}

	// 查找数据
	newTask := task.NewTask(jpm.URL, jpm.KeyWords)
	if mongodb.Exists(newTask) {
		mongodb.FindOne(newTask)
	} else {
		// 获取资源出错
		logger.Error.Println("[GET] 获取资源出错")
		rt = ReturnType{
			Status:  jobNotExists,
			Message: fmt.Sprintf("未找到`job_id=%s`的任务", j),
			Result:  ""}
		c.JSON(http.StatusBadRequest, rt)
		mongodb.DeleteOne(j)
		return
	}

	if newTask.Status == util.TaskCompleted {
		// 任务已完成
		r := resource.Resource{URL: newTask.URL}
		mongodb.FindOne(&r)

		at := abstext.NewAbsText(newTask.URL, newTask.KeyWords)
		mongodb.FindOne(at)

		av := absvideo.AbsVideo{URL: newTask.URL}
		mongodb.FindOne(&av)

		prefix := r.Location
		pictures := make(map[string]string)

		for _, filename := range av.Abstract {
			pictures[filename] = ImagesToBase64(prefix + filename)
		}

		rt = ReturnType{
			Status:  int(newTask.Status),
			Message: util.GetTaskStatus(newTask.Status),
			Result: gin.H{
				"text":     at.Abstract,
				"pictures": pictures,
			}}
		c.JSON(http.StatusOK, rt)
		return
	}

	// 任务已经存在, 且未完成
	rt = ReturnType{
		Status:  int(newTask.Status),
		Message: util.GetTaskStatus(newTask.Status),
		Result:  ""}

	logger.Debug.Printf("[GET] 返回状态{%+v: %+v}.\n", rt.Status, rt.Message)
	c.JSON(http.StatusOK, rt)
}

var GetJobID = func(c *gin.Context) {
	logger.Debug.Println("[GET] 开始")
	var rt ReturnType
	// 获取数据
	jobID := c.Param("job_id")
	newJob := &job.Job{JobID: jobID}
	if mongodb.Exists(newJob) {
		mongodb.FindOne(newJob)
	} else {
		// 获取资源出错
		logger.Error.Println("[GET] 获取资源出错")
		rt = ReturnType{
			Status:  jobNotExists,
			Message: fmt.Sprintf("未找到`job_id=%s`的任务", jobID),
			Result:  ""}
		c.JSON(http.StatusBadRequest, rt)
		return
	}

	newTask := task.NewTask(newJob.URL, newJob.KeyWords)
	if mongodb.Exists(newTask) {
		mongodb.FindOne(newTask)
	} else {
		// 获取资源出错
		logger.Error.Println("[GET] 获取资源出错")
		rt = ReturnType{
			Status:  jobNotExists,
			Message: fmt.Sprintf("未找到`job_id=%s`的任务", jobID),
			Result:  ""}
		c.JSON(http.StatusBadRequest, rt)
		mongodb.DeleteOne(newJob)
		return
	}

	if newTask.Status == util.TaskCompleted {
		// 任务已完成
		r := resource.Resource{URL: newTask.URL}
		mongodb.FindOne(&r)

		at := abstext.NewAbsText(newTask.URL, newTask.KeyWords)
		mongodb.FindOne(at)

		av := absvideo.AbsVideo{URL: newTask.URL}
		mongodb.FindOne(&av)

		prefix := r.Location
		pictures := make(map[string]string)

		for _, filename := range av.Abstract {
			pictures[filename] = ImagesToBase64(prefix + filename)
		}

		rt = ReturnType{
			Status:  int(newTask.Status),
			Message: util.GetTaskStatus(newTask.Status),
			Result: gin.H{
				"text":     at.Abstract,
				"pictures": pictures,
			}}
		c.JSON(http.StatusOK, rt)
		return
	}

	// 任务已经存在, 且未完成
	rt = ReturnType{
		Status:  int(newTask.Status),
		Message: util.GetTaskStatus(newTask.Status),
		Result:  ""}

	logger.Debug.Printf("[GET] 返回状态{%+v: %+v}.\n", rt.Status, rt.Message)
	c.JSON(http.StatusOK, rt)
}

func ImagesToBase64(str_images string) string {
	//读原图片
	ff, _ := os.Open(str_images)
	fileInfo, _ := ff.Stat()
	defer ff.Close()
	sourcebuffer := make([]byte, fileInfo.Size())
	n, _ := ff.Read(sourcebuffer)
	//base64压缩
	sourcestring := base64.StdEncoding.EncodeToString(sourcebuffer[:n])
	return sourcestring
}
