package job_router

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"swc/data/abstext"
	"swc/data/absvideo"
	"swc/data/job"
	"swc/data/resource"
	"swc/data/task"
	"swc/dbs/mongodb"
	"swc/dbs/redis"
	"swc/logger"
	"swc/util"

	"github.com/gin-gonic/gin"
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
			Status:  -1,
			Message: err.Error(),
			Result:  ""}
		c.JSON(http.StatusBadRequest, rt)
		return
	}

	j := job.NewJob(jpm.DeviceID, jpm.URL, jpm.KeyWords)
	if !redis.Exists(j) && !mongodb.Exists(j) {
		// 获取资源出错
		logger.Warning.Println("[GET] 查询了不存在的任务")
		rt = ReturnType{
			Status:  -3,
			Message: fmt.Sprintf("未找到`job_id=%s`的任务", j.JobID),
			Result:  ""}
		c.JSON(http.StatusBadRequest, rt)
		return
	}

	// 查找数据
	newTask := task.NewTask(jpm.URL, jpm.KeyWords)
	err = mongodb.FindOne(newTask)
	if err == nil {
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
			file, _ := os.Open(prefix + filename)
			content, _ := ioutil.ReadAll(file)
			pictures[filename] = string(content)
		}

		rt = ReturnType{
			Status:  int(newTask.Status),
			Message: util.GetJobStatus(newTask.Status),
			Result: gin.H{
				"text":     at.Abstract,
				"pictures": pictures,
			}}
		c.JSON(http.StatusOK, rt)
		return
	}

	if !redis.Exists(newTask) {
		// 获取资源出错
		logger.Warning.Println("[GET] 查询了不存在的任务")
		rt = ReturnType{
			Status:  -3,
			Message: fmt.Sprintf("未找到`job_id=%s`的任务", newTask.TaskID),
			Result:  ""}
		c.JSON(http.StatusBadRequest, rt)
		return
	}
	// 任务已经存在, 且未完成
	redis.FindOne(newTask)
	rt = ReturnType{
		Status:  int(newTask.Status),
		Message: util.GetJobStatus(newTask.Status),
		Result:  ""}

	logger.Debug.Printf("[GET] 返回状态{%+v: %+v}.\n", rt.Status, rt.Message)
	c.JSON(http.StatusOK, rt)
}
