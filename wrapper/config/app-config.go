package config

import (
	"bufio"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

// 存放全部配置项名称的常量以及一些配置初始化逻辑

const (
	// 运行配置前缀
	runPrefix = "run."
	// JavaPath Java运行时路径或者命令
	JavaPath = runPrefix + "java-path"
	// ErrorMessage 错误消息
	ErrorMessage = runPrefix + "error-message"
	// PreParameters 前置运行参数
	PreParameters = runPrefix + "pre-parameters"
	// 日志配置前缀
	logPrefix = "log."
	// WriteLogToFile 是否写入错误日志到文件
	WriteLogToFile = logPrefix + "write-to-file"
	// LogPath 日志文件位置
	LogPath = logPrefix + "path"
)

// ReadYAMLConfig 读取config.yaml中全部配置
//
// path 配置文件所在目录
//
// 当解析配置出错时，返回错误对象
func ReadYAMLConfig(path string) error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	e := viper.ReadInConfig()
	if e != nil {
		return e
	}
	return nil
}

// IsGUIApp 该应用程序是否是GUI应用程序
var IsGUIApp = false

// ReadGUIConfig 读取gui配置文件中内容，并设定给一个全局变量
//
// path 配置文件所在目录
//
// 读取出错时返回错误
func ReadGUIConfig(path string) error {
	file, e1 := os.Open(filepath.Join(path, "gui"))
	if e1 != nil {
		return e1
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	reader := bufio.NewReader(file)
	content, _, e2 := reader.ReadLine()
	if e2 != nil {
		return e2
	}
	if strings.TrimSpace(string(content)) == "y" {
		IsGUIApp = true
	}
	return nil
}