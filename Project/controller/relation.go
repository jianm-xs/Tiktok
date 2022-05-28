package controller

import (
	"Project/dao"
	"Project/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// RelationAction 关注操作
func RelationAction(c *gin.Context) {
	// 获取请求参数
	uid := c.DefaultQuery("user_id", "-1")
	userId, _ := strconv.ParseInt(uid, 10, 64)
	token := c.DefaultQuery("token", "")
	toUserId, _ := strconv.ParseInt(c.DefaultQuery("to_user_id", "-1"), 10, 64)
	actionType, _ := strconv.ParseInt(c.DefaultQuery("action_type", "-1"), 10, 64)
	fmt.Println("=====>", userId, "token : ", token, " toUserId : ", toUserId, " actionType : ", actionType)
	// 参数获取失败
	if userId == -1 || token == "" || toUserId == -1 || actionType == -1 {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: -1,
			StatusMsg:  "failed to obtain parameters!",
		})
		return
	}
	// token 校验失败
	if answer, err := CheckToken(token, uid); err != nil || answer == false {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: -2,
			StatusMsg:  "token error!",
		})
		return
	}
	// 数据库操作
	err := dao.RelationAction(userId, toUserId, actionType)
	if err != nil {
		// 操作失败
		c.JSON(http.StatusOK, models.Response{
			StatusCode: -3,
			StatusMsg:  "update Mysql error!",
		})
		return
	}
	// 操作成功
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
