package dao

import (
	"Project/models"
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
		// 如果是关注操作，插入即可
		favorite := models.Favorite{
			UserID:     userId,
			VideoID:    videoId,
			CreateTime: time.Now(),
		}
		// 插入操作
		err := DB.Debug().Create(&favorite).Error
		return err
	} else {
		err := DB.Debug().Delete(models.Favorite{}, "favorite_id = ? and video_id = ?", userId, videoId).Error
		return err
	}
}
