package util

import (
	"gitee.com/swsk33/jar-to-exe-go-wrapper/config"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// 命令行对象实用方法

// JavaCommand 调用Java的命令（Java命令所在的路径）
var JavaCommand = "java"

// SetupJavaCommand 根据配置，初始化Java命令（Java命令所在的路径）
func SetupJavaCommand() {
	// 如果使用嵌入的JRE，那么直接指定路径为嵌入的JRE路径
	if config.GlobalConfig.Build.UseEmbedJre {
		JavaCommand = filepath.Join(TempDirectory, "jre", "bin", "java")
		return
	}
	// 如果指定了外部便携JRE路径，则转换为绝对路径
	jrePath := config.GlobalConfig.Run.JavaPath
	if jrePath != "java" {
		JavaCommand = filepath.Join(SelfPath, jrePath, "java")
	}
}

// 获取一个通用的cmd命令对象
//
// name 要执行的命令程序或者路径
// arg 命令参数
//
// 返回cmd命令对象
func getCmd(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	// 如果程序为GUI程序，则阻止命令运行时弹出命令行窗口
	if config.GlobalConfig.Build.WinApp {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}
	return cmd
}

// JavaExists 检测java命令是否存在
//
// 返回是否存在
func JavaExists() bool {
	cmd := getCmd(JavaCommand, "-version")
	err := cmd.Run()
	return err == nil
}

// ShowErrorDialog 显示错误弹窗
//
// message 消息
func ShowErrorDialog(message string) {
	messageScript := `
	Add-Type -AssemblyName System.Windows.Forms
	[System.Windows.Forms.MessageBox]::Show("` + message + `", "错误", [System.Windows.Forms.MessageBoxButtons]::OK, [System.Windows.Forms.MessageBoxIcon]::Error)
	`
	command := exec.Command("powershell", "-Command", messageScript)
	_ = command.Run()
}

// GetJarRunCmd 获取一个执行jar文件的命令对象
//
// 返回cmd命令对象和日志文件指针，若未配置输出为日志，则日志文件指针为nil
func GetJarRunCmd() (*exec.Cmd, *os.File) {
	// 解析预先配置参数
	preArgs := config.GlobalConfig.Run.PreParameters
	for i := range preArgs {
		preArgs[i] = strings.TrimSpace(preArgs[i])
	}
	// 组装命令参数为切片
	commandSlice := []string{"-jar", filepath.Join(TempDirectory, "main.jar")}
	commandSlice = append(commandSlice, preArgs...)
	args := os.Args
	length := len(args)
	for i := 1; i < length; i++ {
		commandSlice = append(commandSlice, args[i])
	}
	// 构建cmd对象
	cmd := getCmd(JavaCommand, commandSlice...)
	// 重定向程序输出
	var logFile *os.File = nil
	if config.GlobalConfig.Log.WriteToFile {
		// 如果输出日志，则重定向至文件
		logFile, _ = os.OpenFile(filepath.Join(SelfPath, config.GlobalConfig.Log.Path), os.O_CREATE|os.O_APPEND, 0755)
		// 把标准输出和标准错误重定向至这个文件
		cmd.Stdout = logFile
		cmd.Stderr = logFile
	} else {
		// 如果关闭了日志输出，则把标准输入输出和标准错误全部接入到系统终端里
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	// 无论如何都把标准输入接入到系统终端
	cmd.Stdin = os.Stdin
	return cmd, logFile
}