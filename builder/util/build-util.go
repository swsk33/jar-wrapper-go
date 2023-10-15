package util

import (
	"gitee.com/swsk33/jar-to-exe-go-builder/strategy"
	"github.com/fatih/color"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// 一些构建相关实用函数

// 一些可能需要修改的配置键
const (
	// 构建器相关配置前缀
	buildPrefix = "build."
	// 是否是窗体应用程序
	isGUI = buildPrefix + "win-app"
	// 是否使用嵌入的JRE
	useEmbedJRE = buildPrefix + "use-embed-jre"
)

// BuildIconResource 构建图标资源
//
// iconPath 用户指定的png图标
//
// 构建出错时，返回错误对象
func BuildIconResource(iconPath string) error {
	// 准备资源文件
	winresPath := filepath.Join(WrapperPath, "winres")
	e1 := os.Mkdir(winresPath, 0755)
	if e1 != nil {
		color.Red("创建资源文件夹出错！")
		return e1
	}
	e2 := ExtractEmbedFile("resource/winres-template/user.json", filepath.Join(winresPath, "winres.json"))
	if e2 != nil {
		color.Red("创建清单文件时出错！")
		return e2
	}
	e3 := CopyFile(iconPath, filepath.Join(winresPath, "icon.png"))
	if e3 != nil {
		color.Red("复制资源文件时出错！")
		return e3
	}
	// 开始编译资源
	cmd := exec.Command("go-winres", "make")
	cmd.Dir = WrapperPath
	e4 := cmd.Run()
	if e4 != nil {
		color.Red("编译资源文件时出错！")
		color.Red("可能是图片文件损坏或者超过了256x256大小！")
		return e4
	}
	color.HiYellow("构建资源文件完成！")
	return nil
}

// GenerateJREFolder 根据一个jar文件，调用jlink命令生成一个精简版JRE文件夹
//
// jarPath jar文件路径
// outputJREPath 生成的JRE文件夹路径
//
// 发生错误时返回错误对象
func GenerateJREFolder(jarPath, outputJREPath string) error {
	// 分析jar文件
	color.HiBlue("正在分析jar依赖关系...")
	jdepsCmd := exec.Command("jdeps", "--list-deps", jarPath)
	cmdResult, e1 := jdepsCmd.Output()
	if e1 != nil {
		color.Red("分析jar依赖关系时出错！")
		return e1
	}
	dependencies := strings.Split(string(cmdResult), "\n")
	dependencies = dependencies[:len(dependencies)-1]
	for i := range dependencies {
		dependencies[i] = strings.TrimSpace(dependencies[i])
	}
	// 生成JRE文件夹
	color.HiBlue("正在生成JRE文件夹...")
	jlinkCmd := exec.Command("jlink", "--module-path", filepath.Join("%JAVA_HOME%", "jmods"), "--add-modules", strings.Join(dependencies, ","), "--output", outputJREPath)
	e2 := jlinkCmd.Run()
	if e2 != nil {
		color.Red("生成JRE文件夹时出错！")
		return e2
	}
	color.HiYellow("生成精简JRE完成！")
	return nil
}

// BuildExe 将jar文件构建为exe
//
// gui 是否是窗体应用程序
// arch 构建exe的架构
// jar 原始jar文件路径
// config 指定输入配置文件路径
// output 构建exe的输出位置
// inputEmbedJRE 指定要嵌入的JRE文件夹，如果不使用嵌入的JRE，则该参数传入空字符串""，如果要使用自动嵌入JRE功能，则传入"?"
//
// 构建出错时，返回错误对象
func BuildExe(gui bool, arch, jar, config, output, inputEmbedJRE string) error {
	// 处理路径
	exeOutput := output
	// 如果指定的是相对路径，则转换成绝对路径
	if !filepath.IsAbs(output) {
		var e1 error
		exeOutput, e1 = filepath.Abs(output)
		if e1 != nil {
			color.Red("指定的输出路径有误！")
			return e1
		}
	}
	// 准备文件
	e2 := CopyFile(jar, filepath.Join(WrapperPath, "main.jar"))
	if e2 != nil {
		color.Red("获取jar文件失败！")
		return e2
	}
	e3 := CopyFile(config, filepath.Join(WrapperPath, "config.yaml"))
	if e3 != nil {
		color.Red("获取配置文件失败！")
		return e3
	}
	// 读取config.yaml
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(WrapperPath)
	e4 := viper.ReadInConfig()
	if e4 != nil {
		color.Red("读取运行配置文件失败！")
		return e4
	}
	// 处理嵌入JRE逻辑
	embedJRETargetPath := filepath.Join(WrapperPath, "jre")
	// 如果开启了自动嵌入JRE功能，则生成JRE文件夹并修改配置
	if inputEmbedJRE == "?" {
		e := GenerateJREFolder(jar, embedJRETargetPath)
		if e != nil {
			return e
		}
		viper.Set(useEmbedJRE, true)
	} else {
		// 否则，判断是使用手动指定的JRE进行嵌入还是不使用嵌入JRE
		isEmbedJRE := inputEmbedJRE != ""
		viper.Set(useEmbedJRE, isEmbedJRE)
		// 如果要使用嵌入的JRE，则将嵌入的JRE文件夹也复制到构建目录
		if isEmbedJRE {
			e := CopyFolder(inputEmbedJRE, embedJRETargetPath)
			if e != nil {
				color.Red("复制嵌入JRE文件夹失败！")
				return e
			}
		} else {
			// 否则，生成占位文件
			_ = os.MkdirAll(embedJRETargetPath, 0755)
			file, e1 := os.OpenFile(filepath.Join(embedJRETargetPath, "placeholder"), os.O_CREATE|os.O_WRONLY, 0755)
			if e1 != nil {
				color.Red("创建占位文件失败！")
				return e1
			}
			_, e2 := file.WriteString("1")
			if e2 != nil {
				color.Red("写入占位文件失败！")
				return e2
			}
			_ = file.Close()
		}
	}
	// 获取构建变量
	goArch, e5 := strategy.GetGoArchitecture(arch)
	if e5 != nil {
		return e5
	}
	// 处理GUI应用程序情况
	ldFlags := "-w -s"
	// 修改运行配置中的GUI部分
	viper.Set(isGUI, gui)
	if gui {
		// 添加额外编译符号
		ldFlags += " -H=windowsgui"
	}
	// 构建之前刷入运行配置文件
	e6 := viper.WriteConfig()
	if e6 != nil {
		color.Red("写入运行配置失败！")
		return e6
	}
	// 执行构建命令
	cmd := exec.Command("go", "build", "-ldflags", ldFlags, "-o", exeOutput)
	cmd.Env = append(os.Environ(), goArch)
	cmd.Dir = WrapperPath
	e7 := cmd.Run()
	if e7 != nil {
		color.Red("构建exe时发生错误！")
		return e7
	}
	color.HiYellow("构建exe完成！")
	return nil
}

// CleanTemp 清理构建目录下临时文件
func CleanTemp() {
	_ = os.Remove(filepath.Join(WrapperPath, "main.jar"))
	_ = os.Remove(filepath.Join(WrapperPath, "config.yaml"))
	_ = os.Remove(filepath.Join(WrapperPath, "rsrc_windows_386.syso"))
	_ = os.Remove(filepath.Join(WrapperPath, "rsrc_windows_amd64.syso"))
	_ = os.RemoveAll(filepath.Join(WrapperPath, "winres"))
	_ = os.RemoveAll(filepath.Join(WrapperPath, "jre"))
}