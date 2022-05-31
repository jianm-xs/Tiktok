package controller

import (
	"Project/dao"
	"Project/models"
	"Project/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// RelationAction 关注操作
func RelationAction(c *gin.Context) {
	// 获取请求参数
	token := c.DefaultQuery("token", "")
	toUserId, _ := strconv.ParseInt(c.DefaultQuery("to_user_id", "-1"), 10, 64)
	actionType, _ := strconv.ParseInt(c.DefaultQuery("action_type", "-1"), 10, 64)
	// 参数获取失败
	if token == "" || toUserId == -1 || actionType == -1 {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: -1,
			StatusMsg:  "failed to obtain parameters!",
		})
		return
	}
	var userId int64
	myClaims, err := utils.ParseToken(token)
	if err != nil { // token 解析失败
		c.JSON(http.StatusOK, models.Response{
			StatusCode: -2,
			StatusMsg:  "token error!",
		})
		return
	} else { // 如果 token 解析成功，获取 userId
		userId, _ = strconv.ParseInt(myClaims.Uid, 10, 64)
	}
	// 数据库操作
	err = dao.RelationAction(userId, toUserId, actionType)
	if err != nil {
		// 操作失败
		c.JSON(http.StatusOK, models.Response{
			StatusCode: -3,
			StatusMsg:  err.Error(),
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
	var result models.FollowList // 结果
	var queryId, userId int64
	queryId, _ = strconv.ParseInt(c.Query("user_id"), 10, 64) // 获取请求的 user_id
	token := c.DefaultQuery("token", "")                      // 用户的鉴权 token，可能为空
	myClaims, err := utils.ParseToken(token)
	if err != nil { // token 解析失败
		userId = -1 // 说明 token 无效，设置一个不可能存在的 userID, 这样就不影响查找
	} else { // 如果 token 解析成功，获取 userId
		userId, _ = strconv.ParseInt(myClaims.Uid, 10, 64)
	}
	fmt.Println(userId)
	result.UserList, err = dao.GetFollowList(queryId)
	if err != nil {
		result.Response.StatusCode = -1 // 查询失败
		result.Response.StatusMsg = "search error!"
	}
	result.Response.StatusCode = 0 // 成功，设置状态码和描述
	result.Response.StatusMsg = "success"
	c.JSON(http.StatusOK, result) // 设置返回的信息
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
