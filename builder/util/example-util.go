package util

import "fmt"

// 关于命令示例的实用类

// CommandExample 表示一条命令示例的结构体
type CommandExample struct {
	// 示例完整命令语句
	Command string
	// 示例解释说明
	Note string
}

// ExampleToString 将命令示例转换成字符串，每个命令一行
//
// examples 传入多个示例命令
//
// 返回转换后字符串
func ExampleToString(examples ...*CommandExample) string {
	// 统计命令语句的最大长度
	maxWidth := 0
	for _, item := range examples {
		if len(item.Command) > maxWidth {
			maxWidth = len(item.Command)
		}
	}
	// 生成字符串
	result := ""
	for _, item := range examples {
		result += fmt.Sprintf(fmt.Sprintf("  %%-%ds    %%s\n", maxWidth), item.Command, item.Note)
	}
	return result
}