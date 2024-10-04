package cmd

import (
	"fmt"
	"gitee.com/swsk33/sclog"
	"github.com/spf13/cobra"
)

// 根命令定义
var rootCmd = &cobra.Command{
	Use:   "jar2exe-go",
	Short: "该命令可以快速地将jar文件打包为一个单独的exe文件",
	Long:  "jar2exe-go命令可以快速地将jar文件打包为一个单独的exe文件，并支持通过配置文件灵活定制一些功能",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("欢迎使用jar2exe-go命令！")
		fmt.Println("执行jar2exe-go -h查看帮助")
	},
}

func init() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
}

// ExecuteRoot 执行根命令
func ExecuteRoot() {
	e := rootCmd.Execute()
	if e != nil {
		sclog.ErrorLine("命令执行出错！")
		sclog.Error("%s\n", e)
	}
}