package models

import "time"

// Comment 结构体，定义了评论的基本信息
type Comment struct {
	ID         int64     `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"` // 评论 id
	Author     User      `gorm:"foreignKey:AuthorID;references:ID" json:"user"`  // 评论发布者
	AuthorID   int64     `gorm:"column:author_id" json:"-"`                      // 外键
	VideoID    int64     `gorm:"column:video_id" json:"video_id"`                // 视频 id
	Content    string    `gorm:"column:content" json:"content"`                  // 评论内容
	IsDelete   bool      `gorm:"column:is_delete" json:"is_delete"`              // 是否被删除
	CreateTime time.Time `gorm:"column:create_time" json:"create_date"`          // 创建时间
	UpdateTime time.Time `gorm:"column:update_time" json:"-"`                    // 更新时间
}

// CommentActionResponse 评论操作，返回评论
type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment"` // 评论信息
}

type CommentActionRequest struct {
	UserID      int64  `form:"user_id" binding:"required"`     // 用户id
	Token       string `form:"token" binding:"required"`       // 用户鉴权 token
	VideoID     int64  `form:"video_id" binding:"required"`    // 视频 id
	ActionType  int    `form:"action_type" binding:"required"` // 1-发布评论，2-删除评论
	CommentText string `form:"comment_text"`                   // 用户填写的评论内容，在action_type=1的时候使用
	CommentID   int64  `form:"comment_id"`                     // 要删除的评论id，在action_type=2的时候使用
}
