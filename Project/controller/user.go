// feed 包，该包封装了用户相关的接口
// 创建人：龚江炜
// 创建时间：2022-5-14

package controller

import (
	"Project/models"
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type MyClaims struct {
	Uid string `json:"uid" :"uid"` //用户id
	jwt.StandardClaims
}

//密钥
var MySecret = []byte("密钥")

//生成 Token
func GenToken(uid string) (string, error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		uid,
		jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Hour * 24 * 7)), // 定义过期时间,7天后过期
			Issuer:    "test",                                     // 签发人
		},
	}
	// 对自定义MyClaims加密,jwt.SigningMethodHS256是加密算法得到第二部分
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名,给这个token盐加密,并获得完整的编码后的字符串token
	return token.SignedString(MySecret)
}

type UserLoginResponse struct {
	models.Response        // 状态码、状态描述
	UserId          int64  `json:"user_id,omitempty"` // 用户id
	Token           string `json:"token"`             // 鉴权
}

// Login 用户登录接口
func Login(c *gin.Context) {

	var result UserLoginResponse //登录状态码，默认 0 为成功，其他为失败

	//获取用户名和密码
	userName := c.Query("username")
	userPassword := c.Query("password")

	//数据库连接P
	db, _ := sql.Open("mysql", "root:root@tcp(81.70.17.190:3306)/Tiktok")
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			result.Response.StatusCode = -1 // 更改状态码
			result.Response.StatusMsg = "Close database error!"
			result.UserId = -1
			result.Token = ""
			c.JSON(http.StatusOK, result) // 设置返回的信息
			return
		}
	}(db) // 使用完毕后关闭数据库

	err := db.Ping() //连接数据库

	if err != nil {
		result.Response.StatusCode = -2
		result.Response.StatusMsg = "Connect database error!"
		result.UserId = -1
		result.Token = ""
		c.JSON(http.StatusOK, result)
		return
	}

	//查询数据库username是否存在
	ans, err := db.Query("SELECT name FROM user WHERE name=?;", userName) // 查询user表中username值为用户所输入的
	// 保存查询结果
	if err != nil {
		result.Response.StatusCode = -3
		result.Response.StatusMsg = "Read username error!"
		result.UserId = -1
		result.Token = ""
		c.JSON(http.StatusOK, result)
		return
	}

	var userId int64
	var name string
	var password string

	for ans.Next() {
		err := ans.Scan(&name) //获取查询结果
		if err != nil {
			result.Response.StatusCode = -4
			result.Response.StatusMsg = "Read username error!"
			result.UserId = -1
			result.Token = ""
			c.JSON(http.StatusOK, result)
			return
		}
	}

	//查询数据库中的所有数据
	queryCommand := "SELECT user_id, name, password FROM user where name = ? and password = ?"
	rows, _ := db.Query(queryCommand, userName, userPassword)

	for rows.Next() {
		err := rows.Scan(&userId, &name, &password)
		if err != nil { //如果查询失败
			result.StatusCode = -5
			result.StatusMsg = "query failed!"
			result.UserId = -1
			result.Token = ""
			c.JSON(http.StatusOK, result)
			fmt.Println("失败3")
			return
		}

		tokenStr, _ := GenToken(string(userId))

		//用户名和密码同时匹配成功就返回id和token
		result.StatusCode = 0
		result.StatusMsg = "success"
		result.UserId = userId
		result.Token = tokenStr
		c.JSON(http.StatusOK, result)
		return
	}
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
