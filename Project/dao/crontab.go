// crontab 定义了定时任务的相关定义
// 包含了定时任务的初始化和相关任务
// 创建人：龚江炜
// 创建时间：2022-6-6

package dao

import (
	"Project/models"
	"encoding/json"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm/clause"
	"strconv"
)

// InitTimerTask 初始化定时器并启动
// 返回值：
// 		如果初始化失败，返回错误，否则返回 nil
func InitTimerTask() error {
	// 精确到秒
	crontab := cron.New(cron.WithSeconds())
	spec := "*/30 * * * * ?"
	_, err := crontab.AddFunc(spec, updateAllData)
	if err != nil {
		return err
	}
	crontab.Start()
	return nil
}

// updateAllData 把 Redis 中所有数据写入到 MySQL 中
func updateAllData() {
	// 更新所有数据
	_ = updateFollowData()
	_ = updateFavoriteData()
	_ = updateUserData()
	_ = updateVideoData()
}

// updateFollowData 把 Redis 中关注的所有数据写入到 MySQL 中
// 返回值：
// 		如果操作失败，返回错误，否则返回 nil
func updateFollowData() error {
	data, err := RedisDB.HGetAll("follow").Result()
	if err != nil {
		// 如果获取失败，返回错误
		return err
	}
	for _, val := range data {
		if val != "nil" {
			// 当前关注记录有效，插入到数据库
			// 插入即可
			var follow models.Follow
			err = json.Unmarshal([]byte(val), &follow)
			if err != nil {
				return err
			}
			// 如果存在，更新 create_time， 如果不存在。插入关注记录。按照 ID 查找
			err = DB.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{"create_time"}),
			}).Create(&follow).Error
			if err != nil {
				// 如果插入失败，返回错误
				return err
			}
		}
		// 如果 val == nil 说明当前关注记录已经被删除了，不需要进行操作
	}
	return nil
}

// updateFavoriteData 把 Redis 中点赞的所有数据写入到 MySQL 中
// 返回值：
// 		如果操作失败，返回错误，否则返回 nil
func updateFavoriteData() error {
	data, err := RedisDB.HGetAll("favorite").Result()
	if err != nil {
		// 如果获取失败，返回错误
		return err
	}
	for _, val := range data {
		if val != "nil" {
			// 当前点赞记录有效，插入到数据库
			var favorite models.Favorite
			err = json.Unmarshal([]byte(val), &favorite)
			if err != nil {
				return err
			}
			// 如果存在，更新 create_time， 如果不存在。插入关注记录。按照 ID 查找
			err = DB.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{"create_time"}),
			}).Create(&favorite).Error
			if err != nil { // 数据写入失败
				return err
			}
		}
		// 如果 val == nil 说明当前点赞记录已经被删除了，不需要进行操作
	}
	return nil
}

// updateUserData 把 Redis 中所有用户的数据写入到 MySQL 中
// 返回值：
// 		如果操作失败，返回错误，否则返回 nil
func updateUserData() error {
	// 更新用户的关注数
	data, err := RedisDB.HGetAll("user:followCount").Result()
	if err != nil {
		// 如果获取失败，返回错误
		return err
	}
	var userId, followCount, followerCount int64
	for field, val := range data {
		if val == "nil" {
			// 如果不存在，不管
			continue
		}
		// 获取用户ID
		userId, err = strconv.ParseInt(field, 10, 64)
		if err != nil {
			return err
		}
		// 获取关注数
		followCount, err = strconv.ParseInt(val, 10, 64)
		// 将数据写入数据库
		err = DB.Debug().Model(&models.User{ID: userId}).
			UpdateColumn("follow_count", followCount).Error
		if err != nil {
			return err
		}
	}
	// 更新用户的粉丝数
	data, err = RedisDB.HGetAll("user:followerCount").Result()
	if err != nil {
		// 如果获取失败，返回错误
		return err
	}
	for field, val := range data {
		if val == "nil" {
			// 如果不存在，不管
			continue
		}
		// 获取用户ID
		userId, err = strconv.ParseInt(field, 10, 64)
		if err != nil {
			return err
		}
		// 获取粉丝数
		followerCount, err = strconv.ParseInt(val, 10, 64)
		// 将数据写入数据库
		err = DB.Debug().Model(&models.User{ID: userId}).
			UpdateColumn("follower_count", followerCount).Error
		if err != nil {
			return err
		}
	}
	return nil
}

// updateVideoData 把 Redis 中所有视频的数据写入到 MySQL 中
// 返回值：
// 		如果操作失败，返回错误，否则返回 nil
func updateVideoData() error {
	// 更新视频的点赞数
	data, err := RedisDB.HGetAll("video:favoriteCount").Result()
	if err != nil {
		// 如果获取失败，返回错误
		return err
	}
	var videoId, favoriteCount, commentCount int64
	for field, val := range data {
		// 获取视频ID
		videoId, err = strconv.ParseInt(field, 10, 64)
		if err != nil {
			return err
		}
		// 获取点赞数
		favoriteCount, err = strconv.ParseInt(val, 10, 64)
		// 将数据写入数据库
		err = DB.Debug().Model(&models.Video{ID: videoId}).
			UpdateColumn("favorite_count", favoriteCount).Error
		if err != nil {
			return err
		}
	}
	// 更新用户的粉丝数
	data, err = RedisDB.HGetAll("video:commentCount").Result()
	if err != nil {
		// 如果获取失败，返回错误
		return err
	}
	for field, val := range data {
		// 获取视频 ID
		videoId, err = strconv.ParseInt(field, 10, 64)
		if err != nil {
			return err
		}
		// 获取评论数
		commentCount, err = strconv.ParseInt(val, 10, 64)
		// 将数据写入数据库
		err = DB.Debug().Model(&models.Video{ID: videoId}).
			UpdateColumn("comment_count", commentCount).Error
		if err != nil {
			return err
		}
	}
	return nil
}
