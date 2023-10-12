package main

import (
	"embed"
	"gitee.com/swsk33/jar-to-exe-go-wrapper/config"
	"gitee.com/swsk33/jar-to-exe-go-wrapper/util"
	"github.com/spf13/viper"
	"os"
)

// 嵌入配置文件和jar文件

//go:embed main.jar config.yaml
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
	// 先释放文件至临时文件夹
	e1 := util.ExtractAllFileInEmbedFS(fs, util.TempDirectory)
	if e1 != nil {
		util.ShowErrorDialog("启动失败！" + e1.Error())
		os.Exit(1)
	}
	// 程序运行结束时删除临时文件
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(util.TempDirectory)
	// 初始化Viper加载配置
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(util.TempDirectory)
	e2 := viper.ReadInConfig()
	if e2 != nil {
		util.ShowErrorDialog("启动失败！配置读取错误！")
		os.Exit(1)
	}
	// 检查jre
	if !util.JavaExists() {
		util.ShowErrorDialog(viper.GetString(config.ErrorMessage))
		os.Exit(1)
	}
	// 运行jar文件
	cmd, logFile := util.GetCmd()
	e3 := cmd.Run()
	if e3 != nil {
		util.ShowErrorDialog("运行出现错误！" + e3.Error())
	}
	if logFile != nil {
		_ = logFile.Close()
	}
}