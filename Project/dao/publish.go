// Publish 对视频数据的数据库操作
// 创建人：吴润泽
// 创建时间：2022-5-15

package dao

import (
	"Project/models"
)

func GetPublishListByUserId(userID int64) ([]models.Video, error) {
	var videos []models.Video
	err := DB.Debug().Preload("Author").Where(
		"author_id = ?", userID).Find(&videos).Error
	return videos, err
}

func CreateVideoByUserId(userID int64, data models.Video) error {
	video := models.Video{
		AuthorID: userID,
		PlayUrl:  data.PlayUrl,
		CoverUrl: data.CoverUrl,
	}
	err := DB.Debug().Create(&video).Error
	return err
}
