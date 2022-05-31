package dao

import (
	"Project/models"
	"gorm.io/gorm"
	"time"
)

// FavoriteAction 点赞操作
// 参数 :
//		userId : 请求的用户 id
// 		videoId : 点赞的视频 id
//		actionType : 操作类型： 1 -> 点赞， 2 -> 取消点赞
// 返回值：
//		如果操作成功，返回 nil， 否则返回错误信息
func FavoriteAction(userId, videoId, actionType int64) error {
	if actionType == 1 {
		// 如果是点赞操作，插入即可
		favorite := models.Favorite{
			UserID:     userId,
			VideoID:    videoId,
			CreateTime: time.Now(),
		}
		// 插入操作
		// 点赞数 + 1
		err := DB.Debug().Model(&models.Video{ID: videoId}).UpdateColumn("favorite_count", gorm.Expr("favorite_count + 1")).Error
		if err != nil {
			return err // 视频的点赞数加一
		}
		err = DB.Debug().Create(&favorite).Error // 插入这条点赞记录
		return err
	} else {
		// 视频点赞数 - 1
		err := DB.Debug().Model(&models.Video{ID: videoId}).UpdateColumn("favorite_count", gorm.Expr("favorite_count - 1")).Error
		if err != nil {
			return err // 视频点赞 -1 失败
		}
		err = DB.Debug().Delete(models.Favorite{}, "favorite_id = ? and video_id = ?", userId, videoId).Error // 删除这条点赞记录
		return err
	}
}

// GetFavoriteList 点赞列表
// 参数 :
//		userId : 请求的用户 id
// 返回值：
//		[]models.Video 成功返回点赞列表，失败返回nil
//		error 成功，返回 nil， 否则返回错误信息

func GetFavoriteList(userId int64) []models.Video {
	var videos []models.Video

	// 查找当前用户关注的所有用户
	queryFollow := DB.Select("follow.user_id, 1 as is_follow").
		Where("follower_id = ?", userId).
		Table("follow")

	// 查询点赞
	queryFavorite := DB.Select("video_id,create_time,1 as is_favorite").
		Where("favorite_id = ?", userId).
		Table("favorite")

	DB.Table("video").
		// 预加载 User，给 user 表加上 is_follow 字段再查找
		Preload("Author", func(db *gorm.DB) *gorm.DB {
			return db.Select("user.*, is_follow").
				Joins("LEFT JOIN (?) AS fo ON user.user_id = fo.user_id", queryFollow).
				Table("user")
		}).
		// 联结点赞视频
		Joins("JOIN (?) AS fa ON fa.video_id = video.video_id", queryFavorite).
		// 按照点赞时间降序排列，即时间最晚的在前面
		Order("fa.create_time DESC").
		// 选择返回的字段, video 表中缺少 comment_count 属性，暂时用0替代
		Select("video.*, is_favorite, 0 as comment_count").
		Find(&videos)

	return videos
}
