package models

import "time"

// Comment 结构体，定义了评论的基本信息
type Comment struct {
	ID         int64     `gorm:"column:video_id;primary_key;AUTO_INCREMENT" json:"id"` // 评论 id
	Author     User      `gorm:"foreignKey:AuthorID;references:ID" json:"user"`        // 评论发布者
	AuthorID   int64     `gorm:"column:author_id" json:"-"`                            // 外键
	Content    string    `gorm:"column:content" json:"content"`                        // 评论内容
	CreateTime time.Time `gorm:"column:create_time" json:"create_date"`                // 创建时间
	UpdateTime time.Time `gorm:"-" json:"-"`                                           // 更新时间
}
