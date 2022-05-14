package main

import (
	"Project/controller"
	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	r.Static("/static", "./video") // 设置静态文件路径

	apiRouter := r.Group("/douyin") // 路由分组

	// 基础接口
	apiRouter.GET("/feed/", controller.Feed)                // 视频流接口
	apiRouter.GET("/user/", controller.UserInfo)            //用户信息接口
	apiRouter.GET("/publish/list/", controller.PublishList) //发布列表接口
	apiRouter.POST("/user/register", controller.Register)   // 用户注册接口
	apiRouter.POST("/user/login", controller.Login)         // 用户登录接口
	apiRouter.POST("/publish/action", controller.Publish)   // 投稿接口

}
