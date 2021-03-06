package main

import (
	"swc/router"

	"github.com/gin-gonic/gin"
)

// GinRouter is a router
func GinRouter() (r *gin.Engine) {
	r = gin.Default()

	// 创建任务
	r.POST("/job", router.PostJob)
	return r
}

func main() {
	r := GinRouter()
	// 默认监听本地(ipv4 + ipv6) 8080 端口
	r.Run()
}
