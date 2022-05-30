// common 包，该包定义了响应结构体
// 创建人：龚江炜
// 创建时间：2022-5-19

package models

// Response 响应对象，定义了响应状态的基本内容
type Response struct {
	StatusCode int32  `json:"status_code"` // 响应状态码
	StatusMsg  string `json:"status_msg"`  // 状态描述，可以为空
}

// FeedResponse 视频流接口响应对象，定义了视频流响应的基本内容
type FeedResponse struct {
	Response          // 状态码、状态描述
	NextTime  int64   `json:"next_time"`  // 本次返回的视频中，发布最早的时间。可以为空
	VideoList []Video `json:"video_list"` // 本次返回的视频列表。可以为空
}

// VideoListResponse 发布列表接口响应对象，定义了发布列表响应的基本内容
type VideoListResponse struct {
	Response          // 状态码、状态描述
	VideoList []Video `json:"video_list"` // 本次返回的视频列表。可以为空
}

// UserLoginResponse 用户登录接口响应对象，定义了用户登录响应的基本内容
type UserLoginResponse struct {
	Response        // 状态码、状态描述
	UserId   int64  `json:"user_id,omitempty"` // 用户 id
	Token    string `json:"token"`             // 用户鉴权 token
}

// UserInfoResponse 用户信息接口响应对象，返回状态码和用户信息
type UserInfoResponse struct {
	Response
	User `json:"user"`
}

// CommentListResponse 评论列表响应对象，返回评论列表
type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list"` // 评论列表
}

// CommentActionResponse 评论操作，返回评论
type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment"` // 评论信息
}

// FollowList 关注列表，返回所有关注的用户
type FollowList struct {
	Response
	UserList []User `json:"user_list"` // 关注的用户
}
