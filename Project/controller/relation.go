// relation 包，该包封装了关注相关的接口
// 创建人：龚江炜
// 创建时间：2022-5-25

package controller

import (
	"Project/common"
	"Project/dao"
	"Project/models"
	"Project/utils"
	"github.com/gin-gonic/gin"
	"log"
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
			StatusCode: common.StatusQuery,
			StatusMsg:  "获取参数失败",
		})
		return
	}
	var userId int64
	myClaims, err := utils.ParseToken(token)
	if err != nil { // token 解析失败
		c.JSON(http.StatusOK, models.Response{
			StatusCode: common.StatusToken, // 失败，设置状态码和描述
			StatusMsg:  err.Error(),
		}) // 设置返回的信息
		return
	} else { // 如果 token 解析成功，获取 userId
		userId, _ = strconv.ParseInt(myClaims.Uid, 10, 64)
	}
	// 数据库操作
	err = dao.RelationAction(userId, toUserId, actionType)
	if err != nil {
		// 操作失败
		c.JSON(http.StatusOK, models.Response{
			StatusCode: common.StatusData,
			StatusMsg:  err.Error(),
		})
	} else {
		// 操作成功
		c.JSON(http.StatusOK, models.Response{
			StatusCode: common.StatusOK,
			StatusMsg:  "success!",
		})
	}
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
	result.UserList, err = dao.GetFollowList(queryId, userId)
	if err != nil {
		result.Response.StatusCode = common.StatusData // 查询失败
		result.Response.StatusMsg = err.Error()
	} else {
		result.Response.StatusCode = common.StatusOK // 成功，设置状态码和描述
		result.Response.StatusMsg = "success"
	}
	c.JSON(http.StatusOK, result) // 设置返回的信息
}

// FollowerList 粉丝列表
func FollowerList(c *gin.Context) {
	// 返回的结果
	var result models.FollowList
	var queryId, userId int64
	// 获取请求的 token 和 queryId
	queryId, _ = strconv.ParseInt(c.DefaultQuery("user_id", "-1"), 10, 64)
	token := c.DefaultQuery("token", "")
	// 参数获取失败
	if queryId == -1 {
		result.StatusCode = common.StatusQuery
		result.StatusMsg = "获取参数失败"
		result.UserList = nil
		c.JSON(http.StatusOK, result)
		return
	}

	// 解析 token
	myClaims, err := utils.ParseToken(token)
	// token 解析失败
	if err != nil {
		log.Println(err)
		userId = -1 // 说明 token 无效，设置一个不可能存在的 userID, 这样就不影响查找
	} else { // 如果 token 解析成功，获取 userId
		userId, _ = strconv.ParseInt(myClaims.Uid, 10, 64)
	}

	// 获取粉丝列表
	result.UserList, err = dao.GetFollowerUserList(userId, queryId)
	if err != nil {
		result.StatusCode = common.StatusData
		result.StatusMsg = err.Error()
	} else {
		result.StatusCode = common.StatusOK
		result.StatusMsg = "success!"
	}
	c.JSON(http.StatusOK, result)
}
