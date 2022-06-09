// relation 包，该包包含了关注的数据库操作
// 创建人：龚江炜
// 创建时间：2022-5-25

package dao

import (
	"Project/models"
	"Project/utils"
	"errors"
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
	var err error
	if userId == toUserId { // 如果是自己对自己操作，不管
		return errors.New("you cannot operate on yourself")
	}
	var count int64 // 查看有没有对应的 user-follower 对
	_ = DB.Debug().Table("follow").Where("user_id = ? AND follower_id = ?", toUserId, userId).Count(&count).Error

	userIdStr := strconv.FormatInt(userId, 10)
	toUserIdStr := strconv.FormatInt(toUserId, 10)
	userKey := "user:followCount"
	toUserKey := "user:followerCount"
	RelationKey := userIdStr + ":" + toUserIdStr
	switch actionType {
	case PUBLISH:
		// 如果是关注操作，插入即可
		follow := models.Follow{
			AuthorID:   userId,
			UserID:     toUserId,
			CreateTime: time.Now(),
		}
		// 雪花算法获取 ID
		if follow.ID, err = utils.FollowIDWorker.NextID(); err != nil {
			return err
		}
		// 插入关注记录
		err = CreateData("follow", follow, RelationKey)
		if err != nil {
			// 插入失败，返回错误
			return err
		}
		// 当前用户的关注人数 + 1
		err = IncreaseValue(userKey, models.User{ID: userId}, "follow_count", "user", userIdStr)
		if err != nil {
			// 如果更新数据失败，返回错误
			return err
		}
		// 被关注者的粉丝 + 1
		err = IncreaseValue(toUserKey, models.User{ID: toUserId}, "follower_count", "user", toUserIdStr)
		return err
	case DELETE:
		// 删除关注记录
		err = DeleteData("follow", RelationKey)
		if err != nil {
			// 插入失败，返回错误
			return err
		}
		// 当前用户的关注人数 - 1
		err = DecreaseValue(userKey, models.User{ID: userId}, "follow_count", "user", userIdStr)
		if err != nil {
			// 如果更新数据失败，返回错误
			return err
		}
		// 被关注者的粉丝 - 1
		err = DecreaseValue(toUserKey, models.User{ID: toUserId}, "follower_count", "user", toUserIdStr)
		return err
	default:
		// 防御性
		return errors.New("invalid operation")
	}
}

// GetFollowList 获取关注者列表
// 参数 :
//		queryId : 要查看的用户的 id
//		userId : 当前用户的 id
// 返回值：
//		返回关注者列表和错误信息
func GetFollowList(queryId, userId int64) ([]models.User, error) {
	// 更新数据库数据
	err := updateFollowData()
	if err != nil {
		return nil, err
	}
	var users []models.User // 结果
	// 查找 queryID 关注的用户
	err = DB.Debug().Table("follow").
		Select("user.*").
		Where("follower_id = ?", queryId).
		Joins("LEFT JOIN user ON user.user_id = follow.follow_id").
		Find(&users).Error
	if err != nil {
		// 查询失败，返回错误信息
		return nil, err
	}
	// 更新用户信息
	userIdStr := strconv.FormatInt(userId, 10)
	err = UpdateUsers(users[:], userIdStr)
	return users, err
}

// GetFollowerUserList 查询粉丝信息，返回粉丝列表
// 参数 :
//		userId : 用户的 id
//		queryId : 要查询用户的 id
// 返回值：
//		返回根据 queryId 查询出的粉丝列表
func GetFollowerUserList(userId int64, queryId int64) ([]models.User, error) {
	// 更新数据库数据
	err := updateFollowData()
	if err != nil {
		return nil, err
	}
	var users []models.User // 结果
	// 查找 queryID 的粉丝
	err = DB.Debug().Table("follow").
		Select("user.*").
		Where("follow.user_id = ?", queryId).
		Joins("LEFT JOIN user ON user.user_id = follow.follower_id").
		Find(&users).Error
	// 联结当前用户关注的人，完成 is_follow 查询
	if err != nil {
		// 查询失败
		return nil, err
	}
	// 更新用户信息
	userIdStr := strconv.FormatInt(userId, 10)
	err = UpdateUsers(users[:], userIdStr)
	return users, err
}
