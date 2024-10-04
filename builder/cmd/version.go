package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"runtime"
)

// 输出版本号的子命令
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "输出版本号",
	Long:  "用于输出当前的版本号，x.y.z格式",
	Run: func(cmd *cobra.Command, args []string) {
		color.HiGreen("jar2exe-go version %d.%d.%d", 1, 3, 0)
		color.HiBlue("jar-wrapper-go builder by swsk33")
		color.HiYellow("Build using %s", runtime.Version())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}