package job_router

import (
	"net/http"
	"swc/data/job"
	"swc/data/resource"
	"swc/dbs/mongodb"
	"swc/logger"
	"swc/server/task_builder"
	"sync"

	"github.com/gin-gonic/gin"
)

type ReturnType struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
	Title   string      `json:"title"`
}

// JobPostMessage 用户利用 POST 提交的数据, 用于为任务创建唯一的 ID
type JobPostMessage struct {
	DeviceID string   `json:"device_id"`
	URL      string   `json:"url"`
	KeyWords []string `json:"keywords,omitempty"`
}

var jobMu sync.Mutex

// PostJob 接收并处理发往 /job 的 POST 请求
var PostJob = func(c *gin.Context) {
	logger.Info.Println("[POST Job] 收到 POST 请求")
	var rt ReturnType
	var err error

	// 解析 json 字段
	var jpm JobPostMessage
	err = c.ShouldBindJSON(&jpm)
	if err != nil {
		logger.Error.Printf("[POST Job] 数据解析失败: %+v.\n", err)
		rt = ReturnType{
			Status:  jsonErr,
			Message: err.Error(),
			Result:  ""}
		c.JSON(http.StatusBadRequest, rt)
		return
	}

	logger.Info.Printf("[POST Job] 开始创建任务: %+v.\n", jpm)

	// 构建任务
	newJob := job.NewJob(jpm.DeviceID, jpm.URL, jpm.KeyWords)
	jobMu.Lock()
	logger.Info.Println("[POST Job] newJob", newJob)

	if mongodb.Exists(newJob) {
		jobMu.Unlock()
		// 返回 JobID
		rt = GetRT(0, "任务已存在", gin.H{"job_id": newJob.JobID}, newJob.URL)
		logger.Info.Printf("[POST Job] 任务已存在: %+v.\n", newJob)
	} else {
		err := mongodb.InsertOne(newJob)
		if err != nil {
			logger.Info.Printf("[POST Job] job 插入失败: %+v.\n", err)
		}
		jobMu.Unlock()

		// 开启任务
		go task_builder.AddTask(newJob.URL, newJob.KeyWords)

		// 返回 JobID
		rt = GetRT(0, "任务已接收", gin.H{"job_id": newJob.JobID}, newJob.URL)
		logger.Info.Printf("[POST Job] 任务已接收: %+v.\n", newJob)
	}
	c.JSON(http.StatusOK, rt)
}

func GetRT(status int, m string, r interface{}, url string) (rt ReturnType) {
	rt.Status = status
	rt.Message = m
	rt.Result = r
	rs := resource.Resource{URL: url}
	mongodb.FindOne(&rs)
	if rs.Title == "" {
		rs.Title = "标题获取中..."
	}
	rt.Title = rs.Title
	return
}
