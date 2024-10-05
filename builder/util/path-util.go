package util

import (
	"os"
	"path/filepath"
)

// 一些程序相关的路径常量

// WrapperPath 打包器源代码所在文件夹
var WrapperPath string

// SetupPath 初始化全部路径常量
// 先搜寻自身所在目录下是否存在wrapper文件夹
// 若不存在，则搜索当前运行路径下是否存在wrapper文件夹
func SetupPath() {
	// 寻找自身可执行文件所在文件夹下是wrapper
	path, _ := os.Executable()
	WrapperPath = filepath.Join(filepath.Dir(path), "wrapper")
	if FileExists(WrapperPath) {
		return
	}
	// 否则，寻找当前路径下的
	WrapperPath, _ = filepath.Abs("wrapper")
}