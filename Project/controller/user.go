// feed 包，该包封装了用户相关的接口
// 创建人：龚江炜
// 创建时间：2022-5-14

package controller

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Login 用户登录接口
func Login(c *gin.Context) {

}

// Register 用户注册接口
func Register(c *gin.Context) {

}

// UserInfoResponse 用户信息接口响应对象，返回状态码和用户信息
type UserInfoResponse struct {
	Response
	User `json:"user"`
}

// UserInfo 获取登录用户的 id、昵称等
func UserInfo(c *gin.Context) {
	var result UserInfoResponse    // 结果
	selectId := c.Query("user_id") // 获取请求的 user_id
	// 待用，token 还未完善
	//selectToken := c.DefaultPostForm("token", "-1")
	db, _ := sql.Open("mysql", "root:root@(localhost:3306)/Tiktok") // 设置参数
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			result.Response.StatusCode = -1 // 更改状态码
			result.Response.StatusMsg = "Close database error!"
			c.JSON(http.StatusOK, result) // 设置返回的信息
			return
		}
	}(db) // 使用完毕后关闭数据库
	err := db.Ping() // 连接数据库
	if err != nil {  // 连接失败处理
		result.Response.StatusCode = -2
		result.Response.StatusMsg = "Connect database error!"
		c.JSON(http.StatusOK, result)
		return
	}

	// 以下为数据库连接测试代码，实际功能待实现
	// 预计完善时间：数据库创建完成后完善
	queryCommand := "SELECT user_id, `name` FROM `user` WHERE user_id =" + string(selectId) + ";" // 查询语句
	answer, _ := db.Query(queryCommand)                                                           // 执行查询语句

	for answer.Next() {
		err := answer.Scan(&result.User.Id, &result.User.Name) // 获取查询结果
		if err != nil {                                        // 读取失败处理
			result.Response.StatusCode = -3
			result.Response.StatusMsg = "Read user error!"
			c.JSON(http.StatusOK, result)
			return
		}
	}

	result.Response.StatusCode = 0 // 成功，设置状态码和描述
	result.Response.StatusMsg = "success"
	c.JSON(http.StatusOK, result) // 设置返回的信息
}
