package utils

import "strconv"

// String2int64 将字符串转为 int64,
//							异常或第二返回值全部由占位符接收，慎用
// param:
// 		s: 字符串
// 		foo: 占位
// return: int64
func String2int64(s string, foo ...interface{}) int64 {
	res, _ := strconv.ParseInt(s, 10, 64)
	return res
}
