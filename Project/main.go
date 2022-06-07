// main 包，该包内定义了项目的初始化和程序入口
// 创建人：龚江炜
// 创建时间：2022-5-14

package main

import (
	"Project/config"
	"Project/dao"
	"Project/utils"
	"github.com/gin-gonic/gin"
)

func init() {
	// 连接数据库
	err := dao.InitMysql()
	if err != nil {
		panic(err)
	}
	// 连接 Redis 服务器
	err = dao.InitRedis()
	if err != nil {
		panic(err)
	}
	// 初始化所有 ID 生成器
	err = utils.InitIDWorker()
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
