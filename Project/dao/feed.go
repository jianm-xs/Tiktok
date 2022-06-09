// feed 包，该包包含了视频流相关的数据库操作
// 创建人：龚江炜
// 创建时间：2022-5-15

package dao

import (
	"Project/models"
	"strconv"
	"time"
)

// GetVideos : 执行数据库查询，查找时间小于 lastTime 的前 30 个视频
// 参数 :
//      lastTime : 视频最晚时间，可以为空
//		userId : 查询的用户，用于查询是否关注和点赞，可以为空
// 返回值：
//		查询出来的视频列表和错误信息
func GetVideos(lastTime string, userId int64) ([]models.Video, error) {
	var videos []models.Video
	if len(lastTime) == 0 { // 如果时间为空，获取当前时间
		lastTime = time.Now().Format("2006-01-02 15:04:05")
	}

	// 查询视频列表和作者信息
	err := DB.Debug().Table("video").Limit(30).
		// 预加载 User，给 user 表加上 is_follow 字段再查找
		Preload("Author").
		// 按照创建时间降序排列，即时间最晚的在前面
		Order("video.create_time DESC").
		// 筛选条件，lastTime 之前的视频
		Where("video.create_time < ? ", lastTime).
		Find(&videos).Error
	if err != nil { // 查询失败
		return nil, err
	}
	// 使用 Redis 中的数据更新视频信息
	userIdStr := strconv.FormatInt(userId, 10)
	err = UpdateVideos(videos[:], userIdStr)
	if err != nil {
		// 如果更新出现问题，返回错误
		return nil, err
	}
	return videos, nil
}
