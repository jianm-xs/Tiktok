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
