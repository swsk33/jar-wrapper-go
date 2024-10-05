package config

import (
	"github.com/spf13/viper"
)

// 存放全部配置项名称的常量以及一些配置初始化逻辑

// AppConfig 配置对象结构体
type AppConfig struct {
	// 运行相关配置
	Run struct {
		// Java的运行路径或者命令
		JavaPath string `mapstructure:"java-path"`
		// 没有检测到Java运行环境时的提示内容
		ErrorMessage string `mapstructure:"error-message"`
		// 前置运行参数
		PreParameters []string `mapstructure:"pre-parameters"`
	} `mapstructure:"run"`
	// 日志相关配置
	Log struct {
		// 是否写入日志到文件
		WriteToFile bool `mapstructure:"write-to-file"`
		// 日志输出路径
		Path string `mapstructure:"path"`
	} `mapstructure:"log"`
	// 构建器相关选项，由构建器自动修改
	Build struct {
		// 是否为窗体应用程序
		WinApp bool `mapstructure:"win-app"`
		// 是否使用内嵌JRE
		UseEmbedJre bool `mapstructure:"use-embed-jre"`
	} `mapstructure:"build"`
}

// GlobalConfig 全局配置对象
var GlobalConfig AppConfig

// ReadYAMLConfig 读取config.yaml中全部配置
//
// path 配置文件所在目录
//
// 当解析配置出错时，返回错误对象
func ReadYAMLConfig(path string) error {
	// 设定配置路径
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	// 读取配置
	e := viper.ReadInConfig()
	if e != nil {
		return e
	}
	// 绑定至结构体
	e = viper.Unmarshal(&GlobalConfig)
	if e != nil {
		return e
	}
	return nil
}