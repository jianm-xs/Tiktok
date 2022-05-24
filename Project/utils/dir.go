// utils 包，工具函数
// 创建人：吴润泽
// 创建时间：2022-5-24

package utils

import (
	"os"
	"path/filepath"
	"time"
)

// MkDailyDir 创建一个当天的目录
// 参数：
//    basePath: 起始路径
// 返回值：
//    路径形如 basePath/yyyy/mm/dd/
func MkDailyDir(basePath string) string {
	folderName := time.Now().Format("2006/01/02")
	folderPath := filepath.Join(basePath, folderName)
	_ = os.MkdirAll(folderPath, os.ModePerm)
	return folderPath
}
