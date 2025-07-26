# Docker 使用说明

本文档介绍如何使用 Docker 运行 Google Cloud Storage 上传工具。

## 前提条件

- 安装 [Docker](https://docs.docker.com/get-docker/)
- 安装 [Docker Compose](https://docs.docker.com/compose/install/)
- 准备好 Google Cloud 服务账号凭证文件

## 构建 Docker 镜像

```bash
docker-compose build
```

## 使用 Docker 运行工具

### 设置环境变量

```bash
# 设置存储桶名称
export BUCKET_NAME=your-bucket-name

# 设置凭证文件路径
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/credentials.json
```

### 上传单个文件

```bash
docker-compose run --rm gsupload file /data/your-file.txt
```

### 上传整个目录

```bash
docker-compose run --rm gsupload dir /data/your-directory
```

### 列出存储桶中的对象

```bash
docker-compose run --rm gsupload list
```

### 获取对象的元数据

```bash
docker-compose run --rm gsupload info object-name
```

### 设置对象的过期时间

```bash
docker-compose run --rm gsupload expire object-name 7
```

### 删除对象

```bash
docker-compose run --rm gsupload delete object-name
```

## 使用自定义命令

您可以通过设置 `COMMAND` 环境变量来运行任何自定义命令：

```bash
export COMMAND="file /data/your-file.txt --prefix=uploads/ --expiration=7"
docker-compose run --rm gsupload
```

## 注意事项

- Docker 容器中的工作目录是 `/data`，它被映射到当前目录
- 凭证文件被映射到容器内的 `/credentials.json`
- 确保您的凭证文件具有正确的权限