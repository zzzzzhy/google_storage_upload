# Google Cloud Storage 上传工具

这是一个用 Go 语言编写的命令行工具，用于将文件或目录上传到 Google Cloud Storage，并支持设置对象过期时间。

## 功能特点

- 上传单个文件到 Google Cloud Storage
- 递归上传整个目录到 Google Cloud Storage
- 支持设置对象过期时间
- 支持自定义对象名称前缀

## 使用前准备

1. 确保已安装 Go 环境（推荐 Go 1.16 或更高版本）
2. 确保已有 Google Cloud 项目和存储桶
3. 设置 Google Cloud 认证：
   - 方法一：设置环境变量 `GOOGLE_APPLICATION_CREDENTIALS` 指向您的服务账号密钥文件
   - 方法二：使用 `--credentials` 参数指定服务账号密钥文件路径

## 安装

```bash
# 克隆仓库
git clone https://github.com/yourusername/google_storage_upload.git
cd google_storage_upload

# 编译
go build -o gsupload cmd/main.go
```

## 使用方法

### 上传单个文件

```bash
./gsupload --bucket=your-bucket-name file /path/to/your/file.txt
```

### 上传整个目录

```bash
./gsupload --bucket=your-bucket-name dir /path/to/your/directory
```

### 设置对象过期时间

```bash
./gsupload --bucket=your-bucket-name --expiration=7 file /path/to/your/file.txt
```

### 设置对象名称前缀

```bash
./gsupload --bucket=your-bucket-name --prefix=uploads/ file /path/to/your/file.txt
```

### 使用自定义凭证文件

```bash
./gsupload --bucket=your-bucket-name --credentials=/path/to/credentials.json file /path/to/your/file.txt
```

## 参数说明

- `--bucket`, `-b`: (必需) Google Cloud Storage 存储桶名称
- `--credentials`, `-c`: (可选) Google Cloud 凭证文件路径
- `--expiration`, `-e`: (可选) 对象过期天数，默认为 0（不过期）
- `--prefix`, `-p`: (可选) 对象名称前缀