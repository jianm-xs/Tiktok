// Publish 对视频数据的数据库操作
// 创建人：吴润泽
// 创建时间：2022-5-15

package dao

import (
	"Project/models"
	"time"
)

// CreateVideoByData 将视频的一些链接存储到数据库
// 参数 :
//	title: 视频标题
//  userId : 作者 id
//	playUrl: 视频播放地址
// 	coverUrl: 视频封面地址
// 返回值：
//		Error(nullable)
func CreateVideoByData(title string, authorId int64, playUrl string, coverUrl string) error {
	// 存储相关路径
	video := models.Video{
		AuthorID:   authorId,   // 作者的 id
		Title:      title,      // 视频标题
		PlayUrl:    playUrl,    // 播放地址
		CoverUrl:   coverUrl,   // 封面地址
		CreateTime: time.Now(), // 获取当前时间插入
	}
	err := DB.Debug().Create(&video).Error
	return err
}
