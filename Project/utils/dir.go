package utils

import (
	"os"
	"path/filepath"
	"time"
)

// MkDailyDir 创建一个当天的目录
func MkDailyDir(basePath string) string {
	folderName := time.Now().Format("2006/01/02")
	folderPath := filepath.Join(basePath, folderName)
	_ = os.MkdirAll(folderPath, os.ModePerm)
	return folderPath
}
