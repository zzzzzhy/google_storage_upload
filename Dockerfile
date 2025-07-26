FROM golang:1.20-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件
COPY go.mod ./
COPY go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /gsupload main.go

# 使用轻量级的 alpine 镜像
FROM alpine:latest

# 安装 ca-certificates，用于 HTTPS 请求
RUN apk --no-cache add ca-certificates

# 从构建阶段复制二进制文件
COPY --from=builder /gsupload /usr/local/bin/gsupload

# 设置工作目录
WORKDIR /data

# 设置入口点
ENTRYPOINT ["gsupload"]

# 默认命令
CMD ["--help"]