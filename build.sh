#!/bin/bash

# 设置编译后的输出目录
OUTPUT_DIR="bin"
mkdir -p $OUTPUT_DIR

# 设置编译环境为 Linux 64位
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

echo "开始编译服务..."

# 编译 log-api
echo "正在编译 log-api..."
go build -o $OUTPUT_DIR/log-api application/log-api/api/log.go

# 编译 log-ingester
echo "正在编译 log-ingester..."
go build -o $OUTPUT_DIR/log-ingester application/log-ingester/rpc/ingester.go

# 编译 log-query
echo "正在编译 log-query..."
go build -o $OUTPUT_DIR/log-query application/log-query/rpc/query.go

# 编译 user-auth
echo "正在编译 user-auth..."
go build -o $OUTPUT_DIR/user-auth application/user-auth/rpc/auth.go

echo "------------------------------------------------"
echo "编译完成！二进制文件位于 $OUTPUT_DIR 目录下。"
echo ""
echo "部署建议："
echo "1. 将 bin/ 目录下的二进制文件上传到 Linux 服务器。"
echo "2. 在服务器上创建一个工作目录，例如 /app。"
echo "3. 将二进制文件放入 /app 中。"
echo "4. 在 /app 下创建 etc/ 目录，并将对应的 .yaml 配置文件放入其中。"
echo "5. 运行命令示例：./log-api -f etc/logapi-api.yaml"
echo "------------------------------------------------"