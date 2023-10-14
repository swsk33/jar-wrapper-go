package config

import (
	"github.com/spf13/viper"
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
	// 构建器生成的配置前缀
	buildPrefix = "build."
	// WinAPP 是否为GUI窗体应用程序
	WinAPP = buildPrefix + "win-app"
	// UseEmbedJRE 是否使用嵌入的JRE
	UseEmbedJRE = buildPrefix + "use-embed-jre"
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