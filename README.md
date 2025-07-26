# Google Cloud Storage 上传工具

这是一个用 Go 语言编写的命令行工具，用于将文件或目录上传到 Google Cloud Storage，并支持设置对象过期时间。

[![发布版本](https://github.com/yourusername/google_storage_upload/actions/workflows/release.yml/badge.svg)](https://github.com/yourusername/google_storage_upload/actions/workflows/release.yml)

## 功能特点

- 上传单个文件到 Google Cloud Storage
- 递归上传整个目录到 Google Cloud Storage（保留目录结构）
- 支持设置对象过期时间
- 支持自定义对象名称前缀
- 列出存储桶中的对象
- 删除存储桶中的对象
- 获取对象的元数据

## 使用前准备

1. 确保已安装 Go 环境（推荐 Go 1.16 或更高版本）
2. 确保已有 Google Cloud 项目和存储桶
3. 设置 Google Cloud 认证：
   - 方法一：设置环境变量 `GOOGLE_APPLICATION_CREDENTIALS` 指向您的服务账号密钥文件
   - 方法二：使用 `--credentials` 参数指定服务账号密钥文件路径

## 安装

### 从发布版本安装

您可以从 [GitHub Releases](https://github.com/yourusername/google_storage_upload/releases) 页面下载预编译的二进制文件。

#### Linux/macOS

```bash
# 下载最新版本（以 Linux amd64 为例）
curl -L https://github.com/yourusername/google_storage_upload/releases/latest/download/gsupload_linux_amd64.tar.gz | tar xz

# 移动到可执行路径
sudo mv gsupload_linux_amd64 /usr/local/bin/gsupload
chmod +x /usr/local/bin/gsupload
```

#### Windows

下载 zip 文件，解压后将可执行文件放在您的 PATH 环境变量包含的目录中。

### 从源代码构建

```bash
# 克隆仓库
git clone https://github.com/yourusername/google_storage_upload.git
cd google_storage_upload

# 编译
go build -o gsupload main.go

# 或者使用 Makefile 构建所有平台
make build-all
```

### 使用 Docker

```bash
# 构建 Docker 镜像
docker-compose build

# 运行
export BUCKET_NAME=your-bucket-name
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/credentials.json
docker-compose run --rm gsupload file /data/your-file.txt
```

## 使用方法

### 上传单个文件

```bash
gsupload --bucket=your-bucket-name file /path/to/your/file.txt
```

### 上传整个目录

```bash
gsupload --bucket=your-bucket-name dir /path/to/your/directory
```

### 设置对象过期时间

```bash
gsupload --bucket=your-bucket-name --expiration=7 file /path/to/your/file.txt
```

### 设置对象名称前缀

```bash
gsupload --bucket=your-bucket-name --prefix=uploads/ file /path/to/your/file.txt
```

### 使用自定义凭证文件

```bash
gsupload --bucket=your-bucket-name --credentials=/path/to/credentials.json file /path/to/your/file.txt
```

### 列出存储桶中的对象

```bash
gsupload --bucket=your-bucket-name list [PREFIX]
```

### 获取对象的元数据

```bash
gsupload --bucket=your-bucket-name info OBJECT_NAME
```

### 设置对象的过期时间

```bash
gsupload --bucket=your-bucket-name expire OBJECT_NAME DAYS
```

### 删除对象

```bash
gsupload --bucket=your-bucket-name delete OBJECT_NAME
```

## 参数说明

- `--bucket`, `-b`: (必需) Google Cloud Storage 存储桶名称
- `--credentials`, `-c`: (可选) Google Cloud 凭证文件路径
- `--expiration`, `-e`: (可选) 对象过期天数，默认为 0（不过期）
- `--prefix`, `-p`: (可选) 对象名称前缀
- `--version`, `-v`: (可选) 显示版本信息

## 发布新版本

要发布新版本，只需创建一个新的 Git 标签并推送到 GitHub：

```bash
git tag v1.0.0
git push origin v1.0.0
```

GitHub Actions 将自动构建程序并将其发布到 GitHub Releases 页面。
