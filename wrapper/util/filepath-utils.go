package util

import (
	"os"
	"path/filepath"
)

// 管理全部文件路径的源文件

// SelfPath 自身可执行文件所在的路径
var SelfPath string

// TempDirectory 缓存文件夹路径
var TempDirectory string

// SetupPath 初始化全部路径
func SetupPath() {
	// 自身可执行文件所在文件夹
	path, _ := os.Executable()
	SelfPath = filepath.Dir(path)
	// 临时文件夹
	TempDirectory = filepath.Join(os.TempDir(), GetUUIDFilename())
}