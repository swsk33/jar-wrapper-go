package main

import (
	"embed"
	"gitee.com/swsk33/jar-to-exe-go-wrapper/config"
	"gitee.com/swsk33/jar-to-exe-go-wrapper/util"
	"github.com/spf13/viper"
	"os"
)

// 嵌入配置文件和jar文件
//
//go:embed main.jar config.yaml gui
var fs embed.FS

func init() {
	// 初始化路径
	util.SetupPath()
	// 创建临时目录
	e := os.MkdirAll(util.TempDirectory, 0755)
	if e != nil {
		os.Exit(1)
	}
}

func main() {
	// 先释放全部文件至临时文件夹
	e1 := util.ExtractAllFileInEmbedFS(fs, util.TempDirectory)
	if e1 != nil {
		util.ShowErrorDialog("启动失败！" + e1.Error())
		os.Exit(1)
	}
	// 程序运行结束时删除临时文件
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(util.TempDirectory)
	// 加载全部配置文件
	e2 := config.ReadYAMLConfig(util.TempDirectory)
	e3 := config.ReadGUIConfig(util.TempDirectory)
	if e2 != nil || e3 != nil {
		util.ShowErrorDialog("启动失败！配置读取错误！")
		os.Exit(1)
	}
	// 检查jre
	if !util.JavaExists() {
		util.ShowErrorDialog(viper.GetString(config.ErrorMessage))
		os.Exit(1)
	}
	// 运行jar文件
	cmd, logFile := util.GetJarRunCmd()
	defer func(logFile *os.File) {
		_ = logFile.Close()
	}(logFile)
	e4 := cmd.Run()
	if e4 != nil {
		util.ShowErrorDialog("运行出现错误！" + e4.Error())
	}
}