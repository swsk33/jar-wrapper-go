package main

import (
	"embed"
	"gitee.com/swsk33/jar-to-exe-go-wrapper/config"
	"gitee.com/swsk33/jar-to-exe-go-wrapper/util"
	"github.com/spf13/viper"
	"os"
)

// 嵌入配置文件、jar文件和便携式jre
//
//go:embed main.jar config.yaml jre
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
	e2 := util.ExtractEmbedFolder(fs, "jre", util.TempDirectory)
	if e1 != nil {
		util.ShowErrorDialog("启动失败！" + e1.Error())
		os.Exit(1)
	}
	if e2 != nil {
		util.ShowErrorDialog("启动失败！" + e2.Error())
		os.Exit(1)
	}
	// 程序运行结束时删除临时文件
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(util.TempDirectory)
	// 加载全部配置文件
	e3 := config.ReadYAMLConfig(util.TempDirectory)
	if e3 != nil {
		util.ShowErrorDialog("启动失败！配置读取错误！" + e3.Error())
		os.Exit(1)
	}
	// 根据配置获得java命令路径
	util.SetupJavaCommand()
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
	e5 := cmd.Run()
	if e5 != nil {
		util.ShowErrorDialog("运行出现错误！" + e5.Error())
	}
}