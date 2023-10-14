#!/bin/bash

# 检查命令行参数
version=$1
if [ -z "$version" ]; then
	echo 请指定版本号！
	exit
fi

# 创建文件夹并清理
mkdir -p ./target
echo 清理上次构建...
rm -rf ./target/*

# 进行资源文件和可执行文件构建
base_name=jar-wrapper-go
cd ./builder
echo 开始构建资源...
go-winres make
echo "开始构建程序(amd64)..."
mkdir -p ../target/amd64
GOOS=windows GOARCH=amd64 go build -ldflags "-w -s" -o "../target/amd64/jar2exe-go.exe"
echo "开始构建程序(i386)..."
mkdir -p ../target/i386
GOOS=windows GOARCH=386 go build -ldflags "-w -s" -o "../target/i386/jar2exe-go.exe"
rm ./*.syso

# 压缩可执行文件
echo 开始压缩可执行文件...
cd ../target
upx -9 ./amd64/jar2exe-go.exe
upx -9 ./i386/jar2exe-go.exe

# 复制包装器源码模板
echo 开始复制包装器模板...
cp -rf ../wrapper/ ./amd64/
cp -rf ../wrapper/ ./i386/
cd ./amd64/wrapper
rm -rf main.jar config.yaml gui *.syso *.md .idea/ winres/ jre/
cd ../../i386/wrapper
rm -rf main.jar config.yaml gui *.syso *.md .idea/ winres/ jre/

# 打包
echo 进行打包...
cd ../../
mkdir -p ./out
7z a -t7z -mx9 "./out/${base_name}-${version}-amd64.7z" ./amd64/*
7z a -t7z -mx9 "./out/${base_name}-${version}-i386.7z" ./i386/*
echo 成功打包至./target/out目录下！
