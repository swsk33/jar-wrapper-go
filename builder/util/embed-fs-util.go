package util

import (
	"bufio"
	"embed"
	"os"
)

// 对于嵌入文件的实用类

// 存放嵌入的文件的对象
var container *embed.FS

// SetEmbedContainer 初始化时调用该函数传入嵌入的文件对象，使得该包下可以访问
//
// fs 传入嵌入的文件对象的指针
func SetEmbedContainer(fs *embed.FS) {
	container = fs
}

// ExtractEmbedFile 从嵌入文件系统中释放文件到指定目录
//
// embedFilePath 嵌入的文件名或者相对路径
// outputPath 释放到路径（完整路径，带文件名）
//
// 读取出错时，返回错误对象
func ExtractEmbedFile(embedFilePath, outputPath string) error {
	// 先读取文件
	content, e1 := container.ReadFile(embedFilePath)
	if e1 != nil {
		return e1
	}
	// 创建输出文件对象
	file, e2 := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY, 0755)
	if e2 != nil {
		return e2
	}
	// 创建写入器
	writer := bufio.NewWriter(file)
	// 写入文件
	_, e3 := writer.Write(content)
	if e3 != nil {
		return e3
	}
	_ = writer.Flush()
	_ = file.Close()
	return nil
}