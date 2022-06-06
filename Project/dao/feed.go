// feed 包，该包包含了视频流相关的数据库操作
// 创建人：龚江炜
// 创建时间：2022-5-15

package dao

import (
	"Project/models"
	"gorm.io/gorm"
	"time"
)

// GetVideos : 执行数据库查询，查找时间小于 lastTime 的前 30 个视频
// 参数 :
//      lastTime : 视频最晚时间，可以为空
//		userId : 查询的用户，用于查询是否关注和点赞，可以为空
// 返回值：
//		查询出来的视频列表，如果查询出错或无视频，返回 nil

func GetVideos(lastTime string, userId int64) []models.Video {
	var videos []models.Video
	if len(lastTime) == 0 { // 如果时间为空，获取当前时间
		lastTime = time.Now().Format("2006-01-02 15:04:05")
	}
	// 查询 follow
	queryFollow := DB.Raw("? UNION ALL ?",
		DB.Select("? as user_id, 1 as is_follow", userId).Table("follow"),                            // 自己不能关注自己
		DB.Select("follow.user_id, 1 as is_follow").Where("follower_id = ?", userId).Table("follow"), // 查找当前用户关注的所有用户
	)
	// 查询评论
	queryComment := DB.Select("video_id, COUNT(1) AS comment_count").Group("video_id").Table("comment")
	// 查询点赞
	queryFavorite := DB.Select("video_id, 1 as is_favorite").Where("favorite_id = ?", userId).Table("favorite")

	DB.Table("video").Limit(30).
		// 预加载 User，给 user 表加上 is_follow 字段再查找
		Preload("Author", func(db *gorm.DB) *gorm.DB {
			return db.Select("user.*, is_follow").
				Joins("LEFT JOIN (?) AS fo ON user.user_id = fo.user_id", queryFollow).
				Table("user")
		}).
		// 选择返回的字段
		Select("video.*, is_favorite, comment_count").
		// 按照创建时间降序排列，即时间最晚的在前面
		Order("video.create_time DESC").
		// 筛选条件，lastTime 之前的视频
		Where("video.create_time < ? ", lastTime).
		// 联结评论数
		Joins("LEFT JOIN (?) AS co ON co.video_id = video.video_id", queryComment).
		// 联结是否点赞
		Joins("LEFT JOIN (?) AS fa ON fa.video_id = video.video_id", queryFavorite).
		Find(&videos)

	if len(videos) == 0 { // 如果没有视频，返回空
		return nil
	}
	return videos
}
