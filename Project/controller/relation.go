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

	// 返回的结果
	var result models.Response

	// 获取请求的 token和userId
	userId := c.DefaultQuery("userId", "")
	token := c.DefaultQuery("token", "")

	//token解析
	myClaims, err := ParseToken(token)
	if err != nil { // token 解析失败
		result.StatusCode = -2            // 失败，设置状态码和描述
		result.StatusMsg = "token error!" // token 有误
		c.JSON(http.StatusOK, result)     // 设置返回的信息
		return
	}

	//返回的粉丝信息
	var userList []models.User
	//如果token解析出的userId与请求中的userId相同，则拥有权限获取粉丝列表
	if myClaims.Uid == userId {
		//获取粉丝列表
		userList = dao.GetFollowerUserList(userId)
	}

	c.JSON(http.StatusOK, models.FollowList{
		Response: models.Response{
			StatusCode: 0,
			StatusMsg:  "success!",
		},
		UserList: userList,
	})
}
