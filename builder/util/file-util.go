package util

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
)

// 文件实用类

// FileExists 判断文件或者文件夹是否存在
//
// filePath 文件或者文件夹路径
func FileExists(filePath string) bool {
	_, e := os.Stat(filePath)
	if e == nil {
		return true
	}
	return !os.IsNotExist(e)
}

// CopyFile 复制文件
//
// origin 原始文件路径
// dest 复制文件目标路径
//
// 若复制过程出现错误，则会返回错误对象
func CopyFile(origin, dest string) error {
	originFile, e := os.OpenFile(origin, os.O_RDONLY, 0755)
	if e != nil {
		return e
	}
	// 先创建文件夹
	_ = os.MkdirAll(filepath.Dir(dest), 0755)
	destFile, e := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0755)
	if e != nil {
		return e
	}
	reader := bufio.NewReader(originFile)
	buffer := make([]byte, 64)
	writer := bufio.NewWriter(destFile)
	for {
		n, e := reader.Read(buffer)
		if e == io.EOF {
			break
		}
		_, _ = writer.Write(buffer[0:n])
		_ = writer.Flush()
	}
	_ = originFile.Close()
	_ = destFile.Close()
	return nil
}

// CopyFolder 复制一整个文件夹
//
// origin 被复制的文件夹路径
// dest 复制文件夹的目标路径
//
// 复制出错时，返回错误对象
func CopyFolder(origin, dest string) error {
	// 打印当前文件夹内文件
	list, e := os.ReadDir(origin)
	if e != nil {
		return e
	}
	// 进行遍历操作
	for _, item := range list {
		// 如果是文件，则复制文件到目的路径
		if !item.IsDir() {
			e := CopyFile(filepath.Join(origin, item.Name()), filepath.Join(dest, item.Name()))
			if e != nil {
				return e
			}
		} else {
			// 如果是文件夹，则进行递归复制
			e := CopyFolder(filepath.Join(origin, item.Name()), filepath.Join(dest, item.Name()))
			if e != nil {
				return e
			}
		}
	}
	return nil
}