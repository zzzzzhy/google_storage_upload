#!/bin/bash

# 安装 Google Cloud Storage 上传工具到 Linux 系统

# 检查是否为 root 用户
if [ "$EUID" -ne 0 ]; then
  echo "请使用 root 权限运行此脚本"
  exit 1
fi

# 设置变量
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/gsupload"
BINARY_NAME="gsupload"
BINARY_PATH="build/gsupload_linux_amd64"

# 检查二进制文件是否存在
if [ ! -f "$BINARY_PATH" ]; then
  echo "找不到 $BINARY_PATH 文件，请先运行 'make build-linux'"
  exit 1
fi

# 创建配置目录
mkdir -p $CONFIG_DIR

# 复制二进制文件到安装目录
cp $BINARY_PATH $INSTALL_DIR/$BINARY_NAME
chmod +x $INSTALL_DIR/$BINARY_NAME

echo "安装完成！"
echo "可执行文件已安装到: $INSTALL_DIR/$BINARY_NAME"
echo "配置目录: $CONFIG_DIR"
echo ""
echo "使用方法:"
echo "  $BINARY_NAME --bucket=your-bucket-name file /path/to/file.txt"
echo "  $BINARY_NAME --bucket=your-bucket-name dir /path/to/directory"
echo ""
echo "如需设置 Google Cloud 凭证，请执行以下命令:"
echo "  export GOOGLE_APPLICATION_CREDENTIALS=\"/path/to/credentials.json\""
echo "或者使用 --credentials 参数指定凭证文件路径"