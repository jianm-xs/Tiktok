// Publish 对视频数据的数据库操作
// 创建人：吴润泽
// 创建时间：2022-5-15

package dao

import (
	"Project/models"
)

// CreateVideoByData 将视频的一些链接存储到数据库
// 参数 :
//	title: 视频标题
//  token: 用来获取用户信息
//	playUrl: 视频播放地址
// 	coverUrl: 视频封面地址
// 返回值：
//		Error(nullable)
func CreateVideoByData(title string, token string, playUrl string, coverUrl string) error {
	// TODO: token
	authorId := 1

	// 存储相关路径
	video := models.Video{
		AuthorID: int64(authorId),
		Title:    title,
		PlayUrl:  playUrl,
		CoverUrl: coverUrl,
	}
	err := DB.Debug().Create(&video).Error
	return err
}
