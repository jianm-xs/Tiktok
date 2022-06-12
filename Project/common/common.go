// package common 项目公用常量、枚举等
// 创建人：吴润泽
// 创建时间：2022-5-23

package common

// 状态相关常量定义

const (
	StatusOK     = 0  // 状态成功
	StatusQuery  = -1 // 获取参数错误
	StatusData   = -2 // 数据操作失败
	StatusToken  = -3 // Token解析失败
	StatusOption = -4 // 无效请求
)
