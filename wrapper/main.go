package main

import (
	"embed"
	"gitee.com/swsk33/jar-to-exe-go-wrapper/config"
	"gitee.com/swsk33/jar-to-exe-go-wrapper/util"
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
		util.ShowErrorDialog("准备临时目录失败！" + e.Error())
		os.Exit(1)
	}
	// 释放全部文件至临时文件夹
	e = util.ExtractAllFileInEmbedFS(fs, util.TempDirectory)
	if e != nil {
		util.ShowErrorDialog("释放运行文件失败！" + e.Error())
		os.Exit(1)
	}
	e = util.ExtractEmbedFolder(fs, "jre", util.TempDirectory)
	if e != nil {
		util.ShowErrorDialog("释放jre失败！" + e.Error())
		os.Exit(1)
	}
	// 读取配置
	e = config.ReadYAMLConfig(util.TempDirectory)
	if e != nil {
		os.Exit(1)
	}
}

func main() {
	// 根据配置获得java命令路径
	util.SetupJavaCommand()
	// 检查jre
	if !util.JavaExists() {
		util.ShowErrorDialog(config.GlobalConfig.Run.ErrorMessage)
		os.Exit(1)
	}
	// 运行jar文件
	cmd, logFile := util.GetJarRunCmd()
	e := cmd.Run()
	if e != nil {
		util.ShowErrorDialog("运行jar出现错误！" + e.Error())
	}
	// 运行完成后，进行清理工作
	if logFile != nil {
		_ = logFile.Close()
	}
	_ = os.RemoveAll(util.TempDirectory)
}