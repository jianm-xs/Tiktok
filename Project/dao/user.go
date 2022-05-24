package dao

import (
	"Project/models"
)

// GetUserInfo 查询用户信息，返回用户的所有公开信息
// 参数 :
//      queryId : 被查询者的 id
//		userId : 用户的 id
// 返回值：
//		查询出来用户，如果查询出错或无该用户，返回 nil

func GetUserInfo(queryId int64, userId int64) models.User {
	var user models.User
	queryFollow := DB.Select("follow.user_id, 1 as is_follow").Where("follower_id = ?", userId).Table("follow")
	DB.Select("user.*, is_follow").
		Where("user.user_id = ?", queryId).
		Joins("LEFT JOIN (?) AS fo ON fo.user_id = user.user_id", queryFollow).
		First(&user)
	return user
}
