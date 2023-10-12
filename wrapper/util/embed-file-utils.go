package util

import (
	"bufio"
	"embed"
	"errors"
	"os"
	"path/filepath"
)

// 关于文件读写的实用类

// ExtractEmbedFile 从嵌入文件系统中释放文件到指定目录
//
// fs 嵌入文件系统对象
// embedFilePath 嵌入的文件名或者相对路径
// outputPath 释放到路径（完整路径，带文件名）
//
// 读取出错时，返回错误对象
func ExtractEmbedFile(fs embed.FS, embedFilePath, outputPath string) error {
	// 先读取文件
	content, e1 := fs.ReadFile(embedFilePath)
	if e1 != nil {
		return errors.New("读取内嵌文件错误！")
	}
	// 创建输出文件对象
	file, e2 := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY, 0755)
	if e2 != nil {
		return errors.New("输出文件至：" + outputPath + "出错，可能没有权限！")
	}
	// 创建写入器
	writer := bufio.NewWriter(file)
	// 写入文件
	_, e3 := writer.Write(content)
	if e3 != nil {
		return errors.New("释放文件到：" + outputPath + "时失败！可能没有权限！")
	}
	_ = writer.Flush()
	_ = file.Close()
	return nil
}

// ExtractAllFileInEmbedFS 释放一个嵌入文件系统中全部文件，嵌入文件夹会被忽略
//
// fs 嵌入文件系统对象
// outputFolder 输出文件夹
//
// 释放出错时返回错误对象
func ExtractAllFileInEmbedFS(fs embed.FS, outputFolder string) error {
	items, _ := fs.ReadDir(".")
	for _, item := range items {
		if item.IsDir() {
			continue
		}
		e := ExtractEmbedFile(fs, item.Name(), filepath.Join(outputFolder, item.Name()))
		if e != nil {
			return e
		}
	}
	return nil
}