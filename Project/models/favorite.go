package models

import "time"

// Favorite 结构体，封装了关注的相关属性
type Favorite struct {
	ID         int64     `gorm:"column:id;primary_key;AUTO_INCREMENT"`  // 关注 id
	User       User      `gorm:"foreignKey:UserID;references:ID"`       // 粉丝信息
	UserID     int64     `gorm:"column:favorite_id"`                    // 外键，粉丝 id
	Video      User      `gorm:"foreignKey:VideoID;references:ID"`      // 被关注者信息
	VideoID    int64     `gorm:"column:video_id"`                       // 外键，被关注者 id
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"` // 创建时间
}
