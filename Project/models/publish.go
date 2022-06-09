// publish 包，该包定义了与 `upload` 表对应的结构体以及视频相关接口请求响应格式
// 创建人：吴润泽
// 创建时间：2022-5-15

package models

import (
	"time"
)

// Video 视频对象，定义了视频的基本信息
type Video struct {
	ID            int64     `gorm:"column:video_id;primary_key;" json:"id"`          // 视频 id
	Author        User      `gorm:"foreignKey:AuthorID;references:ID" json:"author"` // 视频发布者
	AuthorID      int64     `gorm:"column:author_id" json:"-"`                       // 外键
	PlayUrl       string    `gorm:"column:play_url" json:"play_url"`                 // 视频播放地址
	CoverUrl      string    `gorm:"column:cover_url" json:"cover_url"`               // 视频封面地址
	FavoriteCount int64     `gorm:"column:favorite_count" json:"favorite_count"`     // 视频点赞总数
	CommentCount  int64     `gorm:"column:comment_count" json:"comment_count"`       // 视频评论总数
	IsFavorite    bool      `gorm:"-" json:"is_favorite"`                            // 是否已点赞
	Title         string    `gorm:"column:title" json:"title"`                       // 视频标题
	CreateTime    time.Time `gorm:"column:create_time" json:"create_time"`           // 创建时间
	UpdateTime    time.Time `gorm:"-" json:"-"`                                      // 更新时间
}

// PublishVideoRequest 投稿视频请求格式
type PublishVideoRequest struct {
	Title string `json:"title"`
	Token string `json:"token"`
	Data  Video  `json:"data"` // 视频数据
}
