// common 包，该包包含了格式转换等函数
// 创建人：龚江炜
// 创建时间：2022-5-25

package utils

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
)

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

type JsonTime time.Time // 用于 Json中时间格式化

// MarshalJSON 实现它的json序列化方法
func (this JsonTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(this).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

// Value 在写入数据库的时候调用，将 JsonTime 格式改为时间戳格式
func (t JsonTime) Value() (driver.Value, error) {
	tlt := time.Time(t)
	return tlt, nil
}

// Scan 在数据查询出来之前对数据进行相关操作，转换为 JsonTime
func (t *JsonTime) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		*t = JsonTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
