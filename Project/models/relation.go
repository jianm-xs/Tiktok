// relation 包，该包包含了关注的相关模型定义
// 创建人：龚江炜
// 创建时间：2022-5-25

package models

import "time"

// Follow 结构体，封装了关注的相关属性
type Follow struct {
	ID         int64     `gorm:"column:id;primary_key;AUTO_INCREMENT"`  // 关注 id
	Author     User      `gorm:"foreignKey:AuthorID;references:ID"`     // 粉丝信息
	AuthorID   int64     `gorm:"column:follower_id"`                    // 外键，粉丝 id
	User       User      `gorm:"foreignKey:UserID;references:ID"`       // 被关注者信息
	UserID     int64     `gorm:"column:user_id"`                        // 外键，被关注者 id
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"` // 创建时间
}
