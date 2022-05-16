// Package common 定义models层共用的一些结构体
// 创建人：吴润泽
// 创建时间：2022-5-15

package models

// Response 响应对象，定义了响应的基本内容
type Response struct {
	StatusCode int32       `json:"StatusCode"`          // 响应状态码
	StatusMsg  string      `json:"StatusMsg,omitempty"` // 状态描述，可以为空
	ErrorMsg   string      `json:"ErrorMsg,omitempty"`  // 错误信息，可以为空
	Data       interface{} `json:"Data,omitempty"`      // 数据
}
