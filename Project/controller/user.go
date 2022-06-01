// feed 包，该包封装了用户相关的接口
// 创建人：龚江炜
// 创建时间：2022-5-14

package controller

import (
	"Project/dao"
	"Project/models"
	"Project/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// Login : 用户登录接口，返回 用户 id 和 token
// 参数 :
//      c : 返回的信息（状态、用户 id 和 token）

func Login(c *gin.Context) {
	var result models.UserLoginResponse //登录状态码，默认 0 为成功，其他为失败

	//获取用户名和密码
	userName := c.Query("username")
	userPassword := c.Query("password")
	// 数据库查询，判断用户名和密码是否匹配
	// 如果匹配，返回 user_id 到 result 中
	err := dao.UserLogin(&result.UserId, userName, userPassword)
	if err != nil {
		// 用户名和密码匹配不成功，则返回错误信息
		result.StatusCode = -1
		result.StatusMsg = "username or password error!"
		c.JSON(http.StatusOK, result)
		return
	}
	// 根据 user_id 获取 token 返回
	result.Token, err = utils.GenToken(strconv.FormatInt(result.UserId, 10))
	if err != nil {
		// 生成 token 失败，则返回错误信息
		result.StatusCode = -2
		result.StatusMsg = "Get token error!"
		c.JSON(http.StatusOK, result)
		return
	}
	//用户名和密码同时匹配成功，且 token 生成成功。返回 id 和 token
	result.StatusCode = 0
	result.StatusMsg = "success"
	c.JSON(http.StatusOK, result)
}

// Register 用户注册接口
func Register(c *gin.Context) {
	var result models.UserLoginResponse //结果
	username := c.Query("username")
	password := c.Query("password")
	// 查询数据库中 username 是否存在，不存在创建新用户返回 id
	uid := dao.UserRegister(username, password)
	if uid == -1 { // 注册失败
		result.Response.StatusCode = -1
		result.Response.StatusMsg = "fail to register！"
		result.UserId = uid
		result.Token = ""
		c.JSON(http.StatusOK, result)
		return
	}
	// 生成token
	tokenStr, _ := utils.GenToken(strconv.FormatInt(uid, 10))

	result.Response.StatusCode = 0
	result.Response.StatusMsg = "success！"
	result.UserId = uid
	result.Token = tokenStr
	c.JSON(http.StatusOK, result)
	return
}

// UserInfo 获取登录用户的 id、昵称等
func UserInfo(c *gin.Context) {
	var result models.UserInfoResponse // 结果
	var queryId, userId int64
	queryId, _ = strconv.ParseInt(c.Query("user_id"), 10, 64) // 获取请求的 user_id
	token := c.DefaultQuery("token", "")                      // 用户的鉴权 token，可能为空
	myClaims, err := utils.ParseToken(token)
	if err != nil { // token 解析失败
		userId = -1 // 说明 token 无效，设置一个不可能存在的 userID, 这样就不影响查找
	} else { // 如果 token 解析成功，获取 userId
		userId, _ = strconv.ParseInt(myClaims.Uid, 10, 64)
	}
	// 数据库查询结果，获取用户信息
	result.User = dao.GetUserInfo(queryId, userId)
	result.Response.StatusCode = 0 // 成功，设置状态码和描述
	result.Response.StatusMsg = "success"
	c.JSON(http.StatusOK, result) // 设置返回的信息
}
