// user 包，该包定义了与 `users` 表对应的结构体以及用户相关接口请求响应格式
// 创建人：吴润泽
// 创建时间：2022-5-15

package models

import (
	"time"
)

// User 用户对象，定义了用户的基本信息
type User struct {
	ID            int64     `gorm:"column:user_id;primary_key" json:"id"`        // 视频 id
	Name          string    `gorm:"column:name" json:"name"`                     // 用户名
	Password      string    `gorm:"column:password" json:"-"`                    // 密码
	FollowCount   int64     `gorm:"column:follow_count" json:"follow_count"`     // 用户的关注总数
	FollowerCount int64     `gorm:"column:follower_count" json:"follower_count"` // 用户的粉丝总数
	IsFollow      bool      `gorm:"-" json:"is_follow"`                          // 当前登录用户是否关注了该用户
	CreateTime    time.Time `gorm:"column:create_time" json:"create_time"`       // 创建时间
	UpdateTime    time.Time `gorm:"column:update_time" json:"update_time"`       // 更新时间
}
