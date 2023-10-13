package util

import (
	"os"
	"path/filepath"
)

// 一些程序相关的路径常量

// SelfPath 可执行文件自身所在文件夹
var SelfPath string

// WrapperPath 打包器源代码所在文件夹
var WrapperPath string

// SetupPath 初始化全部路径常量
func SetupPath() {
	// 自身可执行文件所在文件夹
	path, _ := os.Executable()
	SelfPath = filepath.Dir(path)
	// 打包器源码路径
	WrapperPath = filepath.Join(SelfPath, "wrapper")
}