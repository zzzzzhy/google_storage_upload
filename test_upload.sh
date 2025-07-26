#!/bin/bash

# 设置变量
BUCKET_NAME="your-bucket-name"  # 替换为您的存储桶名称
TEST_FILE="test_file.txt"
TEST_DIR="test_dir"

# 创建测试文件和目录
echo "创建测试文件和目录..."
echo "这是一个测试文件" > $TEST_FILE
mkdir -p $TEST_DIR
echo "这是测试目录中的文件1" > $TEST_DIR/file1.txt
echo "这是测试目录中的文件2" > $TEST_DIR/file2.txt
mkdir -p $TEST_DIR/subdir
echo "这是子目录中的文件" > $TEST_DIR/subdir/file3.txt

# 上传单个文件
echo "上传单个文件..."
./gsupload --bucket=$BUCKET_NAME file $TEST_FILE

# 上传整个目录（不指定前缀，将使用目录名称作为前缀）
echo "上传整个目录（使用目录名称作为前缀）..."
./gsupload --bucket=$BUCKET_NAME dir $TEST_DIR

# 上传整个目录（指定前缀）
echo "上传整个目录（使用自定义前缀）..."
./gsupload --bucket=$BUCKET_NAME --prefix=custom-prefix/ dir $TEST_DIR

# 列出对象
echo "列出存储桶中的对象..."
./gsupload --bucket=$BUCKET_NAME list

# 列出特定前缀的对象
echo "列出特定前缀的对象..."
./gsupload --bucket=$BUCKET_NAME list test_dir/

echo "测试完成！"
