// favorite 包，该包包含了点赞的相关模型定义
// 创建人：龚江炜
// 创建时间：2022-5-25

package models

import "time"

// Favorite 结构体，封装了点赞的相关属性
type Favorite struct {
	ID         int64     `gorm:"column:id;primary_key;AUTO_INCREMENT"`  // 关注 id
	User       User      `gorm:"foreignKey:UserID;references:ID"`       // 粉丝信息
	UserID     int64     `gorm:"column:favorite_id"`                    // 外键，粉丝 id
	Video      User      `gorm:"foreignKey:VideoID;references:ID"`      // 被关注者信息
	VideoID    int64     `gorm:"column:video_id"`                       // 外键，被关注者 id
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"` // 创建时间
}
