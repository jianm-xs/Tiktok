package dao

import (
	"Project/models"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// UserLogin 用户登录，查询数据库中的信息是否匹配
// 参数 :
//      userId : 返回的结果
//		username : 请求的用户名
// 		userPassword : 请求的用户密码
// 返回值：
//		如果请求的用户名和密码正确，返回 nil ，否则返回 错误

func UserLogin(userId *int64, username string, userPassword string) error {
	var user models.User
	// 指定只查找 user_id 和 password，其他数据不需要
	DB.Select([]string{"user_id", "password"}).
		// 筛选 name = username 的数据
		Where(&models.User{Name: username}).
		First(&user)
	// 判断请求的密码是否和数据库中的密码相匹配
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userPassword))
	if err != nil {
		// 如果密码不匹配，返回 错误
		log.Println(err)
		return err
	}
	// 设置返回的 user_id
	*userId = user.ID
	return nil

}

// GetUserInfo 查询用户信息，返回用户的所有公开信息
// 参数 :
//      queryId : 被查询者的 id
//		userId : 用户的 id
// 返回值：
//		查询出来用户，如果查询出错或无该用户，返回 nil

func GetUserInfo(queryId int64, userId int64) models.User {
	var user models.User // 结果
	// 查看是否关注该用户
	queryFollow := DB.Select("follow.user_id, 1 as is_follow").Where("follower_id = ?", userId).Table("follow")
	// 查询该用户信息
	DB.Select("user.*, is_follow").
		// 条件筛选，按 user_id 查找
		Where("user.user_id = ?", queryId).
		// 联结是否关注
		Joins("LEFT JOIN (?) AS fo ON fo.user_id = user.user_id", queryFollow).
		First(&user)
	return user
}
