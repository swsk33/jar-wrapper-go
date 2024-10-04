package main

import (
	"embed"
	"gitee.com/swsk33/jar-to-exe-go-builder/cmd"
	"gitee.com/swsk33/jar-to-exe-go-builder/util"
)

// configTemplate 存放嵌入的配置文件的对象
//
//go:embed resource/config.yaml resource/winres-template/user.json
var configTemplate embed.FS

func init() {
	util.SetupPath()
	util.SetEmbedContainer(&configTemplate)
}

func main() {
	cmd.ExecuteRoot()
}