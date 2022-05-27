package controller

import (
	"Project/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RelationAction 关注操作
func RelationAction(c *gin.Context) {
	c.JSON(http.StatusOK, models.Response{
		StatusCode: 0,
		StatusMsg:  "success!",
	})
}

// FollowList 关注列表
func FollowList(c *gin.Context) {
	c.JSON(http.StatusOK, models.FollowList{
		Response: models.Response{
			StatusCode: 0,
			StatusMsg:  "success!",
		},
		UserList: nil,
	})
}

// FollowerList 粉丝列表
func FollowerList(c *gin.Context) {
	c.JSON(http.StatusOK, models.FollowList{
		Response: models.Response{
			StatusCode: 0,
			StatusMsg:  "success!",
		},
		UserList: nil,
	})
}
