package dao

import (
	"Project/models"
	"time"
)

// RelationAction 关注操作
// 参数 :
//		userId : 请求的用户 id
// 		toUserId : 关注的用户 id
//		actionType : 操作类型： 1 -> 关注， 2 -> 取消关注
// 返回值：
//		如果操作成功，返回 nil， 否则返回错误信息
func RelationAction(userId, toUserId, actionType int64) error {
	if actionType == 1 {
		// 如果是关注操作，插入即可
		follow := models.Follow{
			AuthorID:   userId,
			UserID:     toUserId,
			CreateTime: time.Now(),
		}
		// 插入操作
		err := DB.Debug().Create(&follow).Error
		return err
	} else {
		err := DB.Debug().Delete(models.Follow{}, "user_id = ? and follower_id = ?", toUserId, userId).Error
		return err
	}
}

// GetFollowList 获取关注者列表
// 参数 :
//		userId : 请求的用户 id
// 		toUserId : 关注的用户 id
//		actionType : 操作类型： 1 -> 关注， 2 -> 取消关注
// 返回值：
//		如果操作成功，返回 nil， 否则返回错误信息
func GetFollowList(userId int64) ([]models.User, error) {
	var users []models.User // 结果
	// 查询该用户信息
	err := DB.Debug().Table("follow").
		Select("user.*").
		// 条件筛选，按 user_id 查找
		Where("follow.follower_id = ?", userId).
		// 联结用户表
		Joins("LEFT JOIN user ON user.user_id = follow.user_id").
		Find(&users).Error
	return users, err
}
