// feed 包，该包封装了用户相关的接口
// 创建人：龚江炜
// 创建时间：2022-5-14

package controller

import (
	"Project/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
)

//**************************************************jwt使用***************************************************
// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 这里额外记录username、password两字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中

type MyClaims struct {
	Uid string `json:"uid"` //用户id
	jwt.StandardClaims
}

// 密钥

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

// 解析 Token

func ParseToken(tokenStr string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	if err != nil {
		fmt.Println(" token parse err:", err)
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

//校验 token

func CheckToken(tokenStr string, uid string) (bool, error) {
	//解析Token
	claims, err := ParseToken(tokenStr)

	// 可加相对应的验证逻辑
	if err == nil {
		//用户请求中的uid与解析token出来的uid进行比较，相等就验证成功
		if uid == claims.Uid {
			return true, nil
		}
	}
	return false, err
}

// 刷新 Token

func RefreshToken(tokenStr string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}

	token, err := jwt.ParseWithClaims(tokenStr, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = jwt.At(time.Now().Add(time.Minute * 10))
		return GenToken(claims.Uid)
	}
	return "", errors.New("couldn't handle this token")
}

//*****************************************************************************************************

// Login 用户登录接口
func Login(c *gin.Context) {
	var result models.UserLoginResponse //登录状态码，默认 0 为成功，其他为失败

	//获取用户名和密码
	userName := c.Query("username")
	userPassword := c.Query("password")

	//数据库连接
	db, _ := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/Tiktok")
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

		tokenStr, _ := GenToken(strconv.FormatInt(userId, 10))

		//用户名和密码同时匹配成功就返回id和token
		result.StatusCode = 0
		result.StatusMsg = "success"
		result.UserId = userId
		result.Token = tokenStr
		c.JSON(http.StatusOK, result)
		return
	}

	//匹配不成功，则返回错误信息
	result.StatusCode = -5
	result.StatusMsg = "userName is error!"
	result.UserId = -1
	result.Token = ""
	c.JSON(http.StatusOK, result)
}

// Register 用户注册接口
func Register(c *gin.Context) {
	var result models.UserLoginResponse //结果
	username := c.Query("username")
	password := c.Query("password")
	fmt.Println("username:", username, "   password:", password)
	db, _ := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/Tiktok") // 设置参数
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
	err := db.Ping() // 连接数据库
	if err != nil {  // 连接失败处理
		result.Response.StatusCode = -2
		result.Response.StatusMsg = "Connect database error!"
		result.UserId = -1
		result.Token = ""
		c.JSON(http.StatusOK, result)
		return
	}

	//查询数据库username是否存在
	ans, err := db.Query("SELECT name FROM user WHERE name=?;", username) // 查询user表中username值为用户所输入的
	var name string = ""                                                  // 保存查询结果
	if err != nil {
		fmt.Println("error:", err)
	}
	for ans.Next() {
		//获取查询结果
		err := ans.Scan(&name)
		if err != nil {
			result.Response.StatusCode = -3
			result.Response.StatusMsg = "Read username error!"
			result.UserId = -1
			result.Token = ""
			c.JSON(http.StatusOK, result)
			return
		}
	}
	//username已经存在
	if name == username {
		result.Response.StatusCode = -4
		result.Response.StatusMsg = "username had exist!"
		result.UserId = -1
		result.Token = ""
		c.JSON(http.StatusOK, result)
		return
	}
	//使用bcrypt对密码加密
	pd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	//在数据库中创建用户数据
	ret, err := db.Exec("insert into user (name,password,follow_count,follower_count,create_time,update_time) values(?,?,0,0,?,?)", username, pd, time.Now(), time.Now())
	if err != nil {
		fmt.Println("error:", err)
		result.Response.StatusCode = -5
		result.Response.StatusMsg = "create user error!"
		result.UserId = -1
		result.Token = ""
		c.JSON(http.StatusOK, result)
		return
	}
	uid, err := ret.LastInsertId()       // 获取注册用户的uid
	tokenStr, _ := GenToken(string(uid)) // 生成token

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
	selectId := c.Query("user_id")     // 获取请求的 user_id
	// 待用，token 还未完善
	//selectToken := c.DefaultPostForm("token", "-1")
	db, _ := sql.Open("mysql", "root:root@(127.0.0.1:3306)/Tiktok") // 设置参数
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
		err := answer.Scan(&result.User.ID, &result.User.Name) // 获取查询结果
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
