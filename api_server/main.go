package main

import (
	"swc/router"
	"swc/util"

	"github.com/gin-gonic/gin"
)

// GinRouter is a router
func GinRouter() (r *gin.Engine) {
	util.SetWorkSpace()
	r = gin.Default()

	// 创建任务
	r.POST("/job", router.PostJob)

	// 查询任务
	r.GET("/job", router.GetJob)

	// 删除任务
	r.DELETE("/job", router.DeleteJob)

	return r
}

func main() {
	r := GinRouter()
	// 默认监听本地(ipv4 + ipv6) 8080 端口
	r.Run()
}
