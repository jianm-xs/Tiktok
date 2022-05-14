// common 包，该包封装了项目共用的一些结构体
// 创建人：龚江炜
// 创建时间：2022-5-14

package controller

// Response 响应对象，定义了响应的基本内容
type Response struct {
	StatusCode int32  `json:"status_code"`          // 响应状态码
	StatusMsg  string `json:"status_msg,omitempty"` // 状态描述，可以为空
}

// Video 视频对象，定义了视频的基本信息
type Video struct {
	Id            int64  `json:"id"`             // 视频 ID
	Author        User   `json:"author"`         // 视频发布者
	PlayUrl       string `json:"play_url"`       // 视频播放地址
	CoverUrl      string `json:"cover_url"`      // 视频封面地址
	FavoriteCount int64  `json:"favorite_count"` // 视频点赞总数
	CommentCount  int64  `json:"comment_count"`  // 视频评论总数
	IsFavorite    bool   `json:"is_favorite"`    // 是否已点赞
}

// User 用户对象，定义了用户的基本信息
type User struct {
	Id            int64  `json:"id"`            // 用户 ID
	Name          string `json:"name"`          // 用户名
	FollowCount   int64  `json:"follow_count"`  // 用户的关注总数
	FollowerCount int64  `json:"follower_cont"` // 用户的粉丝总数
	IsFollow      bool   `json:"is_follow"`     // 当前登录用户是否关注了该用户
}
