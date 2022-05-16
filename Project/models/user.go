// 定义与 `users` 表对应的结构体以及用户相关接口请求响应格式
// 创建人：吴润泽
// 创建时间：2022-5-15

package models

import "gorm.io/gorm"

// User 用户对象，定义了用户的基本信息
type User struct {
	gorm.Model    `gorm:"embedded"`
	Name          string `gorm:"column:name" json:"Name"`                    // 用户名
	FollowCount   int64  `gorm:"column:follow_count" json:"FollowCount"`     // 用户的关注总数
	FollowerCount int64  `gorm:"column:follower_count" json:"FollowerCount"` // 用户的粉丝总数
	IsFollow      bool   `gorm:"column:is_follow" json:"IsFollow"`           // 当前登录用户是否关注了该用户
}
