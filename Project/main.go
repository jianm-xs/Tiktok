package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default() // 创建 Gin 引擎
	initRouter(r)
	err := r.Run(":1080") // 启动 Gin 引擎，监听 1080 端口
	if err != nil {
		return
	}
}
