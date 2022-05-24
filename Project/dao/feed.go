package dao

import (
	"Project/models"
	"gorm.io/gorm/clause"
	"time"
)

// GetVideos : 执行数据库查询，查找时间小于 lastTime 的前 30 个视频
// 参数 :
//      lastTime : 视频最晚时间，可以为空
//		userId : 查询的用户，用于查询是否关注和点赞，可以为空
// 返回值：
//		查询结果，如果查询出错，返回 nil

func GetVideos(lastTime string, userId int64) []models.Video {
	var videos []models.Video
	if len(lastTime) == 0 { // 如果时间为空，获取当前时间
		lastTime = time.Now().Format("2006-01-02 15:04:05")
	}
	queryFollow := DB.Select("follow.user_id, 1 as is_follow").Where("follower_id = ?", userId).Table("follow")
	queryFavorite := DB.Select("video_id, 1 as is_favorite").Where("favorite_id = ?", userId).Table("favorite")

	DB.Debug().Table("video").Limit(30).
		Preload(clause.Associations).
		Select("video.*, user.*, fo.is_follow, fa.is_favorite").
		Order("video.create_time DESC").
		Where("video.create_time < ? ", lastTime).
		Joins("LEFT JOIN (?) AS fa ON fa.video_id = video.video_id", queryFavorite).
		Joins("LEFT JOIN user ON user.user_id = video.author_id").
		Joins("LEFT JOIN (?) AS fo ON fo.user_id = user.user_id", queryFollow).
		Find(&videos)
	for i := 0; i < len(videos); i++ {
		videos[i].Author.IsFollow = videos[i].ISFollow // 赋值
	}
	if len(videos) == 0 { // 如果没有视频，返回空
		return nil
	}
	return videos
}
