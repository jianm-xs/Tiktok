// router 包，该包定义了进行了路由初始化和静态文件路径
// 创建人：龚江炜
// 创建时间：2022-5-14

package main

import (
	"Project/config"
	"Project/controller"
	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
	"time"
)

func initRouter(r *gin.Engine) {
	redisStore := persist.NewRedisStore(redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    config.RedisCfg.Addr,
	}))

	r.StaticFS("upload", http.Dir("./upload")) // 设置静态文件路径
	apiRouter := r.Group("/douyin")            // 路由分组

	// 基础接口
	apiRouter.GET("/feed/",
		cache.CacheByRequestURI(redisStore, 5*time.Second),
		controller.Feed,
	) // 视频流接口
	apiRouter.GET("/user/", controller.UserInfo) //用户信息接口
	apiRouter.GET("/publish/list/",
		cache.CacheByRequestURI(redisStore, 5*time.Second),
		controller.PublishList,
	) // 发布列表接口
	apiRouter.POST("/user/register/", controller.Register) // 用户注册接口
	apiRouter.POST("/user/login/", controller.Login)       // 用户登录接口
	apiRouter.POST("/publish/action/", controller.Publish) // 投稿接口

	// 拓展接口 I
	apiRouter.GET("/favorite/list/", controller.FavoriteList)      // 点赞列表
	apiRouter.GET("/comment/list/", controller.CommentList)        // 评论列表
	apiRouter.POST("/comment/action/", controller.CommentAction)   // 评论操作
	apiRouter.POST("/favorite/action/", controller.FavoriteAction) // 点赞操作

	// 拓展接口 II
	apiRouter.GET("/relation/follow/list/", controller.FollowList)     //关注列表
	apiRouter.GET("/relation/follower/list/", controller.FollowerList) // 粉丝列表
	apiRouter.POST("/relation/action/", controller.RelationAction)     // 关注操作

}
