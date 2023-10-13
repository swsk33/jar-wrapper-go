package util

import (
	"errors"
	"strings"
)

// 命令行参数解析实用类

// GetParameterIndex 返回参数在参数列表中的下标，从0开始
//
// arg 参数
// args 全部参数列表
//
// 返回arg在切片args中的下标，若arg不存在于args中返回-1
func GetParameterIndex(arg string, args []string) int {
	// 处理空格
	arg = strings.TrimSpace(arg)
	for i := range args {
		if arg == strings.TrimSpace(args[i]) {
			return i
		}
	}
	return -1
}

// IsParameterExists 判断参数列表中某个参数是否存在
//
// arg 要判断的参数
// args 传入的全部参数列表
//
// 返回arg是否存在于args中
func IsParameterExists(arg string, args []string) bool {
	return GetParameterIndex(arg, args) != -1
}

// GetParameterNext 获取一个参数后面传入的值
//
// argPrefix 前置参数
// args 传入的全部参数列表
//
// 例如：GetParameterNext("-c", []string{"build", "-c", "./config", "-n", "1"})，返回"./config"
//
// 返回argPrefix的后一个参数的参数值，如果argPrefix不存在或者其后没有值，则返回错误
func GetParameterNext(argPrefix string, args []string) (string, error) {
	// 若参数不存在
	if !IsParameterExists(argPrefix, args) {
		return "", errors.New("缺少参数：" + argPrefix)
	}
	// 定位参数
	index := GetParameterIndex(argPrefix, args)
	// 如果参数后未接值
	if index == len(args)-1 {
		return "", errors.New("参数：" + argPrefix + " 后缺少值！")
	}
	return args[index+1], nil
}