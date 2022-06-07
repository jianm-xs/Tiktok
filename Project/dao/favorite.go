// favorite 包，该包包含了点赞的数据库操作
// 创建人：龚江炜
// 创建时间：2022-5-15

package dao

import (
	"Project/models"
	"errors"
	"gorm.io/gorm"
	"strconv"
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
	var count int64 // 查看有没有对应的 video-user 对
	DB.Table("favorite").Where("favorite_id = ? AND video_id = ?", userId, videoId).Count(&count)
	if count+actionType == ActionError {
		// count 只有 1 和 0
		// 如果 count = 0,那么不能删除(actionType != 2)
		// 如果 count = 1，那么不可以继续插入(actionType != 1)
		return errors.New("action error")
	}
	// 通过 VideoId 获取对应的 Key 值，方便 Redis 操作
	key := "video_favoriteCount_" + strconv.FormatInt(videoId, 10)
	switch actionType {
	case PUBLISH:
		// 点赞操作，插入这条记录
		favorite := models.Favorite{
			UserID:     userId,
			VideoID:    videoId,
			CreateTime: time.Now(),
		}
		// 插入这条点赞记录到数据库中
		err := DB.Debug().Create(&favorite).Error
		if err != nil { // 数据写入失败
			return err
		}
		// 该视频的点赞数 + 1
		err = IncreaseValue(key, models.Video{ID: videoId}, "favorite_count", "video")
		return err
	case DELETE:
		// 删除操作，删除这条记录
		err := DB.Debug().Delete(models.Favorite{}, "favorite_id = ? and video_id = ?", userId, videoId).Error // 删除这条点赞记录
		if err != nil {                                                                                        // 删除点赞记录失败
			return err
		}
		// 该视频的点赞数 - 1
		err = DecreaseValue(key, models.Video{ID: videoId}, "favorite_count", "video")
		return err
	default:
		// 防御性
		return errors.New("invalid operation")
	}
}

// GetFavoriteList 点赞列表
// 参数 :
//		userId : 请求的用户 id
// 返回值：
//		[]models.Video 成功返回点赞列表，失败返回nil
//		error 成功，返回 nil， 否则返回错误信息
func GetFavoriteList(authorId, userId int64) ([]models.Video, error) {
	var videos []models.Video

	// 查询 follow
	queryFollow := DB.Raw("? UNION ALL ?",
		DB.Select("? as user_id, 1 as is_follow", userId).Table("follow"),                            // 自己不能关注自己
		DB.Select("follow.user_id, 1 as is_follow").Where("follower_id = ?", userId).Table("follow"), // 查找当前用户关注的所有用户
	)

	// 查询点赞
	queryFavorite := DB.Select("video_id,create_time,1 as is_favorite").
		Where("favorite_id = ?", authorId).
		Table("favorite")

	err := DB.Debug().Table("video").
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
		Find(&videos).Error
	if err != nil {
		// 如果查询失败，返回错误信息
		return nil, err
	}
	err = UpdateVideos(videos[:])
	if err != nil {
		// 如果更新出现问题，返回错误
		return nil, err
	}

	return videos, nil
}
