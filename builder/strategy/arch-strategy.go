package strategy

import (
	"errors"
	"fmt"
)

// 根据指定的不同架构生成不同的GOARCH环境变量

// 策略容器
var archMap map[string]string

func init() {
	archMap = map[string]string{
		"i386":  "386",
		"amd64": "amd64",
	}
}

// GetGoArchitecture 传入-a参数指定的架构值返回构建时的GOARCH变量值
//
// arch 命令行参数中指定的架构值
//
// 返回完整的GOARCH环境变量语句，若架构不存在返回错误
func GetGoArchitecture(arch string) (string, error) {
	// 获取键值
	value, exist := archMap[arch]
	if !exist {
		return value, errors.New(fmt.Sprintf("指定的架构：%s 不存在！", arch))
	}
	return fmt.Sprintf("GOARCH=%s", value), nil
}