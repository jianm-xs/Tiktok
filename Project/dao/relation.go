package dao

import (
	"Project/models"
	"time"
)

// RelationAction 用户注册，查询用户名是否存在，不存在注册新用户
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
