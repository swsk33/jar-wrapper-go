package cmd

import (
	"errors"
	"gitee.com/swsk33/jar-to-exe-go-builder/util"
	"gitee.com/swsk33/sclog"
	"github.com/spf13/cobra"
	"path/filepath"
)

// 定义标志
var (
	// 配置文件路径
	configFile string
	// exe文件架构
	arch string
	// 图标文件路径
	iconPath string
	// 是否是GUI窗体程序
	gui bool
	// 是否内嵌JRE
	embedJre bool
	// 内嵌JRE路径
	embedJrePath string
	// 是否自动内嵌JRE
	autoEmbedJre bool
)

// 示例命令
var buildExamples = []*util.CommandExample{
	{"jar2exe-go build demo.jar demo.exe", "将当前目录下demo.jar构建打包为demo.exe到当前目录下，使用当前目录下config.yaml作为配置文件"},
	{"jar2exe-go build demo.jar demo.exe -c dir/config.yaml", "将当前目录下demo.jar构建打包为demo.exe，使用./dir/config.yaml作为配置文件"},
	{"jar2exe-go build demo.jar demo.exe -i gopher.png", "将当前目录下demo.jar构建打包为demo.exe，使用gopher.png作为exe图标"},
	{"jar2exe-go build demo.jar demo.exe -g", "将当前目录下demo.jar构建打包为demo.exe，该程序为窗体应用程序"},
	{"jar2exe-go build demo.jar demo.exe -a i386", "将当前目录下demo.jar构建打包为32位的demo.exe可执行文件"},
	{"jar2exe-go build demo.jar demo.exe --embed-jre --embed-jre-path jre17", "将当前目录下demo.jar构建打包为demo.exe，并内嵌JRE文件夹：./jre17作为内嵌的运行环境"},
	{"jar2exe-go build demo.jar demo.exe --auto-embed-jre", "将当前目录下demo.jar构建打包为demo.exe，并自动内嵌JRE运行环境"},
}

// 构建jar为exe的子命令
var buildCmd = &cobra.Command{
	Use:     "build",
	Short:   "构建jar为exe",
	Long:    "将指定的jar文件打包构建为exe文件，命令用法为：\n  jar2exe-go jar文件路径 输出exe路径 [-c 配置文件] [-a 架构] [-i exe图标路径] [-g] [--embed-jre] [--embed-jre-path 指定要嵌入的JRE文件夹] [--auto-embed-jre]",
	Example: util.ExampleToString(buildExamples...),
	// 参数校验
	Args: func(cmd *cobra.Command, args []string) error {
		// 参数个数校验
		if len(args) != 2 {
			sclog.Error("传入参数长度必须为2！当前传入：%d个参数\n", len(args))
			sclog.ErrorLine("第1个参数为打包jar路径")
			sclog.ErrorLine("第2个参数为输出exe路径")
			return errors.New("传入参数长度必须为2！")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// 解析参数
		jarPath := args[0]
		exePath := args[1]
		defer util.CleanTemp()
		// 如果指定了图标，则先构建资源
		sclog.InfoLine("正在构建资源文件...")
		e := util.BuildIconResource(iconPath)
		if e != nil {
			sclog.ErrorLine(e.Error())
			return
		}
		// 判断是否指定了嵌入JRE运行环境
		embedPath := ""
		// 如果指定了使用自动嵌入JRE功能
		if autoEmbedJre {
			embedPath = "?"
		} else {
			// 否则，判断是否是手动指定的要嵌入的JRE
			if embedJre {
				if embedJrePath == "" {
					sclog.ErrorLine(e.Error())
					return
				}
				embedPath = embedJrePath
			}
		}
		// 构建可执行程序
		sclog.InfoLine("正在构建可执行文件...")
		e = util.BuildExe(gui, arch, jarPath, configFile, exePath, embedPath)
		if e != nil {
			sclog.ErrorLine(e.Error())
			return
		}
		sclog.InfoLine("构建完成！")
		outAbsPath, _ := filepath.Abs(exePath)
		sclog.Info("成功生成exe文件至：%s\n", outAbsPath)
		sclog.InfoLine("清理临时文件...")
	},
}

func init() {
	flags := buildCmd.Flags()
	// 绑定标志
	flags.StringVarP(&configFile, "config", "c", "config.yaml", "指定配置文件路径，如果未传入该参数，则默认使用当前路径下的config.yaml作为配置文件")
	flags.StringVarP(&arch, "architecture", "a", "amd64", "指定输出exe文件的架构，可以是以下值：\ni386  输出为32位exe\namd64 输出为64位exe\n当未指定该标志时使用amd64架构")
	flags.StringVarP(&iconPath, "icon", "i", "", "指定exe图标，要求是分辨率不大于256x256的png图片文件")
	flags.BoolVarP(&gui, "gui", "g", false, "当加上该参数时表示被打包的jar是GUI窗体应用程序，那么运行exe时不会显示命令行窗口，默认为命令行程序，运行exe时会显示命令行窗口")
	flags.BoolVarP(&embedJre, "embed-jre", "", false, "当加上该参数时，将Java运行环境(jre)也嵌入至exe中去\n若指定了该标志，则必须指定--embed-jre-path标志的值")
	flags.StringVarP(&embedJrePath, "embed-jre-path", "", "", "指定要嵌入至exe中的JRE文件夹\n如果没有指定--embed-jre标志，那么--embed-jre-path标志无效")
	flags.BoolVarP(&autoEmbedJre, "auto-embed-jre", "", false, "当加上该参数时，构建器将自动分析jar所依赖的模块，并自动生成一个JRE文件夹嵌入到exe中去\n使用该参数时，--embed-jre和--embed-jre-path标志都将无效\n使用该功能要求本地已安装并正确配置了JDK 9及其以上版本的JDK")
	rootCmd.AddCommand(buildCmd)
}