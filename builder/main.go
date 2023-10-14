package main

import (
	"embed"
	"fmt"
	"gitee.com/swsk33/jar-to-exe-go-builder/util"
	"github.com/fatih/color"
	"os"
	"path/filepath"
)

// 打印帮助信息
func printHelp() {
	fmt.Println("\njar2exe-go命令可以快速地将jar文件打包为一个单独的exe文件，并支持通过配置文件灵活定制一些功能")
	fmt.Println("\n用法：")
	fmt.Println("\n* 输出帮助信息：")
	fmt.Println("  jar2exe-go -h")
	fmt.Println("\n* 输出版本信息：")
	fmt.Println("  jar2exe-go -v")
	fmt.Println("\n* 在当前目录下生成一个配置文件模板（config.yaml）：")
	fmt.Println("  jar2exe-go init-config")
	fmt.Println("\n* 打包jar为exe：")
	fmt.Println("  jar2exe-go -j jar文件路径 -o 输出exe路径 [-c 配置文件] [-a 架构] [-i exe图标路径] [-g] [--embed-jre] [--embed-jre-path 指定要嵌入的JRE文件夹]")
	fmt.Println("\n打包jar为exe的参数说明：")
	fmt.Println()
	fmt.Println("  -j : 指定待打包的jar文件路径")
	fmt.Println("  -o : 指定输出的exe文件路径")
	fmt.Println("  -c : 可选参数，指定配置文件路径")
	fmt.Println("       如果未传入该参数，则默认使用当前路径下的config.yaml作为配置文件")
	fmt.Println("  -a : 可选参数，指定输出exe文件的架构，可以是以下值：")
	fmt.Println("       i386  输出为32位exe")
	fmt.Println("       amd64 输出为64位exe")
	fmt.Println("       当未指定-a参数时使用amd64架构")
	fmt.Println("  -i : 可选参数，指定exe图标")
	fmt.Println("       要求是分辨率不大于256x256的png图片文件")
	fmt.Println("  -g : 可选参数，当加上该参数时表示被打包的jar是GUI窗体应用程序，那么运行exe时不会显示命令行窗口")
	fmt.Println("       默认为命令行程序，运行exe时会显示命令行窗口")
	fmt.Println("  --embed-jre      可选参数，当加上该参数时，将Java运行环境(jre)也嵌入至exe中去")
	fmt.Println("                   若指定了该参数，则必须指定--embed-jre-path参数的值")
	fmt.Println("  --embed-jre-path 可选参数，指定要嵌入至exe中的JRE文件夹")
	fmt.Println("                   如果没有指定--embed-jre参数，那么--embed-jre-path参数无效")
}

// configTemplate 存放嵌入的配置文件的对象
//
//go:embed resource/config.yaml resource/winres-template/user.json
var configTemplate embed.FS

func init() {
	util.SetupPath()
	util.SetEmbedContainer(&configTemplate)
}

func main() {
	// 解析命令行参数
	args := os.Args
	// 处理帮助信息
	if util.GetParameterIndex("-h", args) == 1 {
		printHelp()
		return
	}
	// 处理版本信息
	if util.GetParameterIndex("-v", args) == 1 {
		color.HiGreen("jar2exe-go version 1.1.0")
		fmt.Println("jar-wrapper-go builder by swsk33")
		fmt.Println("This builder was developed using Golang")
		return
	}
	// 生成配置模板
	if util.GetParameterIndex("init-config", args) == 1 {
		e := util.ExtractEmbedFile("resource/config.yaml", "./config.yaml")
		if e != nil {
			color.Red("生成配置出错！")
			color.Red(e.Error())
			return
		}
		color.HiGreen("已生成配置文件config.yaml至当前目录下！")
		return
	}
	// 处理构建命令
	jarPath, e1 := util.GetParameterNext("-j", args)
	if e1 != nil {
		color.Red("错误：" + e1.Error())
		printHelp()
		return
	}
	outputPath, e2 := util.GetParameterNext("-o", args)
	if e2 != nil {
		color.Red("错误：" + e2.Error())
		printHelp()
		return
	}
	configPath, e3 := util.GetParameterNext("-c", args)
	if e3 != nil {
		configPath = "./config.yaml"
	}
	archValue, e4 := util.GetParameterNext("-a", args)
	if e4 != nil {
		archValue = "amd64"
	}
	// 如果指定了图标，则先构建资源
	defer util.CleanTemp()
	iconPath, e5 := util.GetParameterNext("-i", args)
	if e5 == nil {
		color.HiBlue("正在构建资源文件...")
		e := util.BuildIconResource(iconPath)
		if e != nil {
			color.Red(e.Error())
			return
		}
	}
	// 判断是否指定了嵌入JRE运行环境
	embedJREPath := ""
	if util.IsParameterExists("--embed-jre", args) {
		path, e := util.GetParameterNext("--embed-jre-path", args)
		if e != nil {
			color.Red(e.Error())
			return
		}
		embedJREPath = path
	}
	// 构建可执行程序
	color.HiBlue("正在构建可执行文件...")
	isGUI := util.IsParameterExists("-g", args)
	e6 := util.BuildExe(isGUI, archValue, jarPath, configPath, outputPath, embedJREPath)
	if e6 != nil {
		color.Red(e6.Error())
		return
	}
	color.HiGreen("构建完成！")
	outAbsPath, _ := filepath.Abs(outputPath)
	color.HiGreen("成功生成exe文件至：%s", outAbsPath)
	color.HiBlue("清理临时文件...")
}