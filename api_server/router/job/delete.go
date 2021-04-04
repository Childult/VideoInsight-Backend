package job_router

import (
	"net/http"
	"swc/data/abstext"
	"swc/data/absvideo"
	"swc/data/job"
	"swc/data/resource"
	"swc/data/task"
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
			Status:  jsonErr,
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
		recursiveDelete(newJob)
		rt = ReturnType{
			Status:  0,
			Message: "删除成功",
			Result:  ""}
	} else {
		rt = ReturnType{
			Status:  jobNotExists,
			Message: "资源不存在",
			Result:  ""}
	}

	// 返回
	logger.Debug.Println("[DELETE] 结束")
	c.JSON(http.StatusBadRequest, rt)
}

func recursiveDelete(j *job.Job) {
	t := task.NewTask(j.URL, j.KeyWords)
	redis.DeleteOne(t)

	r := &resource.Resource{URL: j.URL}
	redis.DeleteOne(r)

	at := abstext.NewAbsText(j.URL, j.KeyWords)
	redis.DeleteOne(at)

	av := &absvideo.AbsVideo{URL: j.URL}
	redis.DeleteOne(av)
}
