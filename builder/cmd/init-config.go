package cmd

import (
	"gitee.com/swsk33/jar-to-exe-go-builder/util"
	"gitee.com/swsk33/sclog"
	"github.com/spf13/cobra"
)

// 初始化配置文件的子命令
var initConfigCmd = &cobra.Command{
	Use:   "init-config",
	Short: "生成配置文件模板",
	Long:  "在当前目录下生成一个配置文件模板，名为config.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		e := util.ExtractEmbedFile("resource/config.yaml", "./config.yaml")
		if e != nil {
			sclog.ErrorLine("生成配置出错！")
			sclog.Error("%s\n", e.Error())
			return
		}
		sclog.InfoLine("已生成配置文件config.yaml至当前目录下！")
	},
}

func init() {
	rootCmd.AddCommand(initConfigCmd)
}