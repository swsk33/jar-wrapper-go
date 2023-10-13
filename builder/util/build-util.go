package util

import (
	"gitee.com/swsk33/jar-to-exe-go-builder/strategy"
	"github.com/fatih/color"
	"os"
	"os/exec"
	"path/filepath"
)

// 一些构建相关实用函数

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

// BuildExe 将jar文件构建为exe
//
// gui 是否是窗体应用程序
// arch 构建exe的架构
// jar 原始jar文件路径
// config 指定配置文件路径
// output 构建exe的输出位置
//
// 构建出错时，返回错误对象
func BuildExe(gui bool, arch, jar, config, output string) error {
	// 处理路径
	commandOutput := output
	// 如果指定的是相对路径，则转换成绝对路径
	if !filepath.IsAbs(output) {
		var e1 error
		commandOutput, e1 = filepath.Abs(output)
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
	e4 := ExtractEmbedFile("resource/gui", filepath.Join(WrapperPath, "gui"))
	if e4 != nil {
		color.Red("准备GUI配置失败！")
		return e4
	}
	// 获取构建变量
	goArch, e5 := strategy.GetGoArchitecture(arch)
	if e5 != nil {
		return e5
	}
	// 处理GUI应用程序情况
	ldFlags := "-w -s"
	if gui {
		// 添加额外编译符号
		ldFlags += " -H=windowsgui"
		// 修改gui配置文件
		file, e1 := os.OpenFile(filepath.Join(WrapperPath, "gui"), os.O_WRONLY, 0755)
		if e1 != nil {
			color.Red("打开GUI配置文件失败！")
			return e1
		}
		_, e2 := file.WriteString("y")
		if e2 != nil {
			color.Red("修改GUI配置文件失败！")
			return e1
		}
		_ = file.Close()
	}
	// 执行构建命令
	cmd := exec.Command("go", "build", "-ldflags", ldFlags, "-o", commandOutput)
	cmd.Env = append(os.Environ(), goArch)
	cmd.Dir = WrapperPath
	e6 := cmd.Run()
	if e6 != nil {
		color.Red("构建exe时发生错误！")
		return e6
	}
	color.HiYellow("构建exe完成！")
	return nil
}

// CleanTemp 清理构建目录下临时文件
func CleanTemp() {
	_ = os.Remove(filepath.Join(WrapperPath, "main.jar"))
	_ = os.Remove(filepath.Join(WrapperPath, "config.yaml"))
	_ = os.Remove(filepath.Join(WrapperPath, "gui"))
	_ = os.Remove(filepath.Join(WrapperPath, "rsrc_windows_386.syso"))
	_ = os.Remove(filepath.Join(WrapperPath, "rsrc_windows_amd64.syso"))
	_ = os.RemoveAll(filepath.Join(WrapperPath, "winres"))
}