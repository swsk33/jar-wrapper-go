package util

import (
	"gitee.com/swsk33/sclog"
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
	if !FileExists(WrapperPath) {
		sclog.ErrorLine("请将包装器代码wrapper文件夹放在可执行文件所在目录下，或者当前运行目录下！")
		sclog.ErrorLine("即将退出...")
		os.Exit(1)
	}
}