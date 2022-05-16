// 定义与 `videos` 表对应的结构体以及视频相关接口请求响应格式
// 创建人：吴润泽
// 创建时间：2022-5-15

package models

import "gorm.io/gorm"

// Video 视频对象，定义了视频的基本信息
type Video struct {
	gorm.Model    `gorm:"embedded"`
	Author        User   `gorm:"foreignKey:AuthorID;references:ID" json:"Author"` // 视频发布者
	AuthorID      int64  `gorm:"column:author_id" json:"AuthorID"`                // 外键
	PlayUrl       string `gorm:"column:play_url" json:"PlayUrl"`                  // 视频播放地址
	CoverUrl      string `gorm:"column:cover_url" json:"CoverUrl"`                // 视频封面地址
	FavoriteCount int64  `gorm:"column:favorite_count" json:"FavoriteCount"`      // 视频点赞总数
	CommentCount  int64  `gorm:"column:comment_count" json:"CommentCount"`        // 视频评论总数
	IsFavorite    bool   `gorm:"column:is_favorite" json:"IsFavorite"`            // 是否已点赞
}

// PublishVideoRequest 投稿视频请求格式
type PublishVideoRequest struct {
	UserID int64  `json:"UserID"`
	Token  string `json:"Token"`
	Data   Video  `json:"Data"` // 视频数据
}

// PublishVideoResponse 投稿视频响应格式
type PublishVideoResponse struct {
}
