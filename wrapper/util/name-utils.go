package util

import (
	"github.com/google/uuid"
	"strings"
)

// 命名实用类

// GetUUIDFilename 生成带有UUID的随机文件名
func GetUUIDFilename() string {
	return "j2e-go-" + strings.Replace(uuid.New().String(), "-", "", -1)
}