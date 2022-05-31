// jwt token 的生成与校验
// 创建人：刘伟欢
// 创建时间：2022-5-30

package utils

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"time"
)

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
// 参数：
//		tokenStr 传入的 token
// 返回：
//		*MyClaims 解析 token 出来的记录的数据
//		err 返回错误

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
