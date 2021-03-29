package main

import (
	"swc/logger"
	job_router "swc/router/job"

	"github.com/gin-gonic/gin"
)

// GinRouter 路由
func GinRouter() (r *gin.Engine) {
	r = gin.Default()

	// 创建任务
	r.POST("/job", job_router.PostJob)

	// 查询任务
	r.GET("/job/:job_id", job_router.GetJob)

	// 删除任务
	r.DELETE("/job", job_router.DeleteJob)

	return r
}

func main() {
	// 初始化日志
	logger.Info.Println("[main] API Server启动")

	r := GinRouter()
	// 默认监听本地(ipv4 + ipv6) 8080 端口
	r.Run()
}
