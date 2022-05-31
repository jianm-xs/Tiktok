package dao

import (
	"Project/models"
	"errors"
	"gorm.io/gorm"
	"strconv"
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
	if userId == toUserId { // 如果是自己对自己操作，不管
		return errors.New("you cannot operate on yourself")
	}
	if actionType == 1 {
		// 如果是关注操作，插入即可
		follow := models.Follow{
			AuthorID:   userId,
			UserID:     toUserId,
			CreateTime: time.Now(),
		}
		// 插入操作
		// 被关注用户的被关注数 + 1
		err := DB.Debug().Model(&models.User{ID: toUserId}).UpdateColumn("follower_count", gorm.Expr("follower_count + 1")).Error
		if err != nil {
			return err
		}
		// 当前用户的关注数 + 1
		err = DB.Debug().Model(&models.User{ID: userId}).UpdateColumn("follow_count", gorm.Expr("follower_count + 1")).Error
		if err != nil {
			return err
		}
		// 插入关注记录
		err = DB.Debug().Create(&follow).Error
		return err
	} else {
		// 被关注用户的被关注数 - 1
		err := DB.Debug().Model(&models.User{ID: toUserId}).UpdateColumn("follower_count", gorm.Expr("follower_count - 1")).Error
		if err != nil {
			return err
		}
		// 当前用户的被关注数 - 1
		err = DB.Debug().Model(&models.User{ID: userId}).UpdateColumn("follow_count", gorm.Expr("follow_count - 1")).Error
		if err != nil {
			return err
		}
		// 删除关注记录
		err = DB.Debug().Delete(models.Follow{}, "user_id = ? and follower_id = ?", toUserId, userId).Error
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
		Select("user.*, 1 as is_follow").
		// 条件筛选，按 user_id 查找
		Where("follow.follower_id = ?", userId).
		// 联结用户表
		Joins("LEFT JOIN user ON user.user_id = follow.user_id").
		Find(&users).Error
	return users, err
}

// GetFollowerUserList 查询粉丝信息，返回粉丝列表
// 参数 :
//		userId : 用户的 id
// 返回值：
//		返回根据用户id查询出的粉丝列表
func GetFollowerUserList(id string) []models.User {
	var user []models.User // 结果
	// 查询用户id
	var userId int64
	userId, _ = strconv.ParseInt(id, 10, 64)
	// 查询 follow
	queryFollow := DB.Raw("? UNION ALL ?",
		DB.Select("? as user_id, 1 as is_follow", userId).Table("follow"),                            // 自己不能关注自己
		DB.Select("follow.user_id, 1 as is_follow").Where("follower_id = ?", userId).Table("follow"), // 查找当前用户关注的所有用户
	)

	// 查询 follower
	DB.Table("user").
		Select("user.*,is_follow").
		// 联结粉丝
		Joins("JOIN `follow` AS f ON user.`user_id`=f.`follower_id` AND f.`user_id`= ? ", userId).
		// 联结关注
		Joins("LEFT JOIN (?) AS fo ON user.user_id = fo.user_id", queryFollow).
		Find(&user)

	return user
}
