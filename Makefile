# Makefile for Google Cloud Storage 上传工具

# 变量
BINARY_NAME=gsupload
GO=go
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date -u '+%Y-%m-%d %H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X 'main.BuildTime=$(BUILD_TIME)'"

# 操作系统和架构
GOOS_LINUX=linux
GOOS_DARWIN=darwin
GOOS_WINDOWS=windows
GOARCH_AMD64=amd64
GOARCH_ARM64=arm64
GOARCH_386=386

# 输出目录
OUTPUT_DIR=build

# 默认目标
.PHONY: all
all: build

# 创建输出目录
$(OUTPUT_DIR):
	mkdir -p $(OUTPUT_DIR)

# 构建当前平台
.PHONY: build
build:
	$(GO) build $(LDFLAGS) -o $(BINARY_NAME) main.go

# 构建 Linux 版本
.PHONY: build-linux
build-linux: $(OUTPUT_DIR)
	GOOS=$(GOOS_LINUX) GOARCH=$(GOARCH_AMD64) $(GO) build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)_linux_amd64 main.go
	GOOS=$(GOOS_LINUX) GOARCH=$(GOARCH_ARM64) $(GO) build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)_linux_arm64 main.go
	GOOS=$(GOOS_LINUX) GOARCH=$(GOARCH_386) $(GO) build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)_linux_386 main.go

# 构建 macOS 版本
.PHONY: build-darwin
build-darwin: $(OUTPUT_DIR)
	GOOS=$(GOOS_DARWIN) GOARCH=$(GOARCH_AMD64) $(GO) build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)_darwin_amd64 main.go
	GOOS=$(GOOS_DARWIN) GOARCH=$(GOARCH_ARM64) $(GO) build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)_darwin_arm64 main.go

# 构建 Windows 版本
.PHONY: build-windows
build-windows: $(OUTPUT_DIR)
	GOOS=$(GOOS_WINDOWS) GOARCH=$(GOARCH_AMD64) $(GO) build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)_windows_amd64.exe main.go
	GOOS=$(GOOS_WINDOWS) GOARCH=$(GOARCH_386) $(GO) build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)_windows_386.exe main.go

# 构建所有平台
.PHONY: build-all
build-all: build-linux build-darwin build-windows

# 安装
.PHONY: install
install:
	$(GO) install $(LDFLAGS)

# 运行测试
.PHONY: test
test:
	$(GO) test ./...

# 清理
.PHONY: clean
clean:
	$(GO) clean
	rm -f $(BINARY_NAME)
	rm -rf $(OUTPUT_DIR)
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
	@echo "  make build         - 构建当前平台的程序"
	@echo "  make build-linux   - 构建 Linux 平台的程序 (amd64, arm64, 386)"
	@echo "  make build-darwin  - 构建 macOS 平台的程序 (amd64, arm64)"
	@echo "  make build-windows - 构建 Windows 平台的程序 (amd64, 386)"
	@echo "  make build-all     - 构建所有平台的程序"
	@echo "  make install       - 安装程序"
	@echo "  make test          - 运行测试"
	@echo "  make clean         - 清理构建文件"
	@echo "  make test-data     - 创建测试数据"
