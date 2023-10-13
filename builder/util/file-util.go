package util

import (
	"bufio"
	"io"
	"os"
)

// 文件实用类

// FileExists 判断文件是否存在
//
// filePath 判断的文件路径
//
// 返回文件是否存在
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
	originFile, e1 := os.OpenFile(origin, os.O_RDONLY, 0755)
	if e1 != nil {
		return e1
	}
	destFile, e2 := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0755)
	if e2 != nil {
		return e2
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