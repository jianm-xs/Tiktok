package main

import (
	"Project/config"
	"Project/dao"
	"Project/utils"
	"github.com/gin-gonic/gin"
)

func init() {
	err := dao.InitMysql()
	if err != nil {
		panic(err)
	}
	// 初始化所有 ID 生成器
	err = utils.InitIdWorker()
	if err != nil {
		panic(err)
	}
}

func main() {
	r := gin.Default() // 创建 Gin 引擎
	initRouter(r)
	err := r.Run(config.ServerPort) // 启动 Gin 引擎，监听 1080 端口
	if err != nil {
		return
	}
}
