#!/bin/fish

# 检查命令行参数
set app_version $argv[1]

if test -z "$app_version"
    echo 请指定版本号！
    exit
end

set app_name jar2exe-go
set package_name jar-wrapper-go

# 创建文件夹
mkdir -p ./target

# 清理上次构建
if test -d ./target/out/
    echo 清理上次构建...
    rm -r ./target/out/
end

# 生成自动补全脚本
echo 生成自动补全脚本...
mkdir -p ./target/completion
cd ./builder
go run main.go completion bash >../target/completion/$app_name-completion.bash
go run main.go completion fish >../target/completion/$app_name.fish

# 进行资源文件和可执行文件构建
echo 开始构建资源...
go-winres make
echo "开始构建程序(amd64)..."
mkdir -p ../target/amd64
GOOS=windows GOARCH=amd64 go build -ldflags "-w -s" -o "../target/amd64/$app_name.exe"
echo "开始构建程序(i386)..."
mkdir -p ../target/i386
GOOS=windows GOARCH=386 go build -ldflags "-w -s" -o "../target/i386/$app_name.exe"
rm ./*.syso

# 压缩可执行文件
echo 开始压缩可执行文件...
cd ../target
upx -9 ./amd64/$app_name.exe
upx -9 ./i386/$app_name.exe

# 复制包装器源码模板
echo 开始复制包装器模板...
cp -rf ../wrapper/ ./amd64/
cp -rf ../wrapper/ ./i386/

# 打包
echo 进行打包...
mkdir -p ./out
7z a -t7z -mx9 -xr!"main.jar" -xr!"config.yaml" -xr!"*.syso" -xr!"*.md" -xr!".idea" -xr!"winres" -xr!"jre" "./out/$package_name-$app_version-amd64.7z" ./amd64/* ./completion/*
7z a -t7z -mx9 -xr!"main.jar" -xr!"config.yaml" -xr!"*.syso" -xr!"*.md" -xr!".idea" -xr!"winres" -xr!"jre" "./out/$package_name-$app_version-i386.7z" ./i386/* ./completion/*
echo 成功打包至./target/out目录下！

# 清理
echo 正在清理...
rm -r ./i386/ ./amd64/ ./completion/
