# Makefile for Google Cloud Storage 上传工具

# 变量
BINARY_NAME=gsupload
GO=go

# 默认目标
.PHONY: all
all: build

# 构建
.PHONY: build
build:
	$(GO) build -o $(BINARY_NAME) main.go

# 安装
.PHONY: install
install:
	$(GO) install

# 运行测试
.PHONY: test
test:
	$(GO) test ./...

# 清理
.PHONY: clean
clean:
	$(GO) clean
	rm -f $(BINARY_NAME)
	rm -rf test_data

# 创建测试数据
.PHONY: test-data
test-data:
	mkdir -p test_data
	echo "这是测试文件1" > test_data/file1.txt
	echo "这是测试文件2" > test_data/file2.txt
	mkdir -p test_data/subdir
	echo "这是子目录中的测试文件" > test_data/subdir/file3.txt

# 帮助
.PHONY: help
help:
	@echo "可用的 make 命令:"
	@echo "  make build    - 构建程序"
	@echo "  make install  - 安装程序"
	@echo "  make test     - 运行测试"
	@echo "  make clean    - 清理构建文件"
	@echo "  make test-data - 创建测试数据"