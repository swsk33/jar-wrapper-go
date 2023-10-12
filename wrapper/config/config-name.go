package config

// 存放全部配置项名称的常量

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