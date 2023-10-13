package util

import (
	"fmt"
	"gitee.com/swsk33/jar-to-exe-go-wrapper/config"
	"github.com/spf13/viper"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// 命令行对象实用方法

// 获取一个通用的cmd命令对象
//
// name 要执行的命令程序或者路径
// arg 命令参数
//
// 返回cmd命令对象
func getCmd(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	// 如果程序为GUI程序，则阻止命令运行时弹出命令行窗口
	if config.IsGUIApp {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}
	return cmd
}

// JavaExists 检测java命令是否存在
//
// 返回是否存在
func JavaExists() bool {
	// 获取配置的Java路径或者命令
	javaCommand := viper.GetString(config.JavaPath)
	var cmd *exec.Cmd = nil
	if javaCommand == "java" {
		cmd = getCmd("java", "-version")
	} else {
		cmd = getCmd(filepath.Join(SelfPath, javaCommand, "java"), "-version")
	}
	err := cmd.Run()
	return err == nil
}

// ShowErrorDialog 显示错误弹窗
//
// message 消息
func ShowErrorDialog(message string) {
	// 准备vbs脚本
	vbsContent := fmt.Sprintf("MsgBox \"%s\", 16, \"错误\"", message)
	// 转换编码
	encoder := simplifiedchinese.GBK.NewEncoder()
	gbkText, _, _ := transform.String(encoder, vbsContent)
	vbsPath := filepath.Join(TempDirectory, "message.vbs")
	vbsFile, _ := os.OpenFile(vbsPath, os.O_WRONLY|os.O_CREATE, 0755)
	_, _ = vbsFile.WriteString(gbkText)
	_ = vbsFile.Close()
	// 执行vbs脚本
	cmd := getCmd("wscript", vbsPath)
	_ = cmd.Run()
}

// GetJarRunCmd 获取一个执行jar文件的命令对象
//
// 返回cmd命令对象和日志文件指针，若未配置输出为日志，则日志文件指针为nil
func GetJarRunCmd() (*exec.Cmd, *os.File) {
	// 获取配置的Java路径或者命令
	javaCommand := viper.GetString(config.JavaPath)
	if javaCommand != "java" {
		javaCommand = filepath.Join(SelfPath, javaCommand, "java")
	}
	// 解析预先配置参数
	preArgs := viper.GetStringSlice(config.PreParameters)
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
	cmd := getCmd(javaCommand, commandSlice...)
	// 重定向程序输出
	writeLog := viper.GetBool(config.WriteLogToFile)
	var logFile *os.File = nil
	if writeLog {
		// 如果输出日志，则重定向至文件
		logFile, _ = os.OpenFile(filepath.Join(SelfPath, viper.GetString(config.LogPath)), os.O_CREATE|os.O_APPEND, 0755)
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