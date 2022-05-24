package dao

import (
	"Project/models"
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
	queryFollow := DB.Select("follow.user_id, 1 as is_follow").Where("follower_id = ?", userId).Table("follow")
	// 查询评论
	queryComment := DB.Select("video_id, COUNT(1) AS comment_count").Group("video_id").Table("comment")
	// follow 和用户关联起来
	queryUser := DB.Select("user.*, is_follow").
		Joins("LEFT JOIN (?) AS fo ON user.user_id = fo.user_id", queryFollow).
		Table("user")
	// 查询点赞
	queryFavorite := DB.Select("video_id, 1 as is_favorite").Where("favorite_id = ?", userId).Table("favorite")

	DB.Table("video").Limit(30).
		Preload("Author").
		Select("video.*, users.*, is_favorite, comment_count").
		Order("video.create_time DESC").
		Where("video.create_time < ? ", lastTime).
		Joins("LEFT JOIN (?) AS co ON co.video_id = video.video_id", queryComment).
		Joins("LEFT JOIN (?) AS fa ON fa.video_id = video.video_id", queryFavorite).
		Joins("LEFT JOIN (?) AS users ON users.user_id = video.author_id", queryUser).
		Find(&videos)

	if len(videos) == 0 { // 如果没有视频，返回空
		return nil
	}
	return videos
}
