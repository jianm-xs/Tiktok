// feed 包，该包封装了用户相关的接口
// 创建人：龚江炜
// 创建时间：2022-5-14

package controller

import (
	"Project/dao"
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
	result.Token, err = GenToken(strconv.FormatInt(result.UserId, 10))
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
	ret, err := db.Exec("insert into user (user_id, name,password,follow_count,follower_count,create_time,update_time) values(3, ?,?,0,0,?,?)", username, pd, time.Now(), time.Now())
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
	var queryId, userId int64
	queryId, _ = strconv.ParseInt(c.Query("user_id"), 10, 64) // 获取请求的 user_id
	token := c.DefaultQuery("token", "")                      // 用户的鉴权 token，可能为空
	myClaims, err := ParseToken(token)
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
