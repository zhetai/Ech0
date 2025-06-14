# # 使用 Node.js Alpine 镜像作为前端构建阶段
# FROM node:22.16.0-alpine AS frontend-build

# # 设置时区为上海
# ENV TZ=Asia/Shanghai

# # 设置前端工作目录
# WORKDIR /app/web

# # 安装PNPM 10.x 版本
# RUN npm install -g pnpm@latest-10

# # 复制前端依赖文件并安装依赖
# COPY ./web/package.json ./web/pnpm-lock.yaml ./
# RUN pnpm install

# # 复制前端源代码并构建项目
# COPY ./web/ .
# RUN pnpm build --mode production 

# # 复制构建后的文件到后端 template 目录
# RUN mv dist /app/template

# # 使用 Golang Alpine 镜像作为后端构建阶段
# FROM golang:1.24.3-alpine AS backend-build

# # 设置后端工作目录
# WORKDIR /app

# # 安装构建时所需的工具（仅限构建阶段）
# # RUN apk add --no-cache gcc musl-dev

# # 设置环境变量启用 cgo
# ENV CGO_ENABLED=0

# # 复制后端代码和必要文件
# COPY ./go.mod ./go.sum ./
# RUN go mod download

# COPY ./cmd ./cmd
# COPY ./internal ./internal
# COPY ./config ./config

# # 创建并设置权限
# RUN mkdir -p /app/data && chmod -R 777 /app/data

# # 构建 Go 后端应用
# RUN go build -o /app/ech0 ./cmd/ech0/main.go

# # 使用更轻量的 Alpine 镜像作为运行时阶段
# FROM alpine:latest AS final

# # 设置时区为上海
# ENV TZ=Asia/Shanghai

# WORKDIR /app

# # ✅ 安装系统根证书，避免 HTTPS 请求失败
# # RUN apk add --no-cache ca-certificates

# # 复制构建阶段的文件
# COPY --from=backend-build /app/config /app/config
# COPY --from=backend-build /app/ech0 /app/ech0
# COPY --from=frontend-build /app/template /app/template

# # 暴露端口
# EXPOSE 6277

# # 运行后端服务
# CMD ["/app/ech0"]

FROM alpine:latest

WORKDIR /app
ENV TZ=Asia/Shanghai

# 获取构建平台信息
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

# 将所有平台的 ech0 二进制和前端资源复制进镜像
COPY /backend-artifacts/* /tmp/
COPY /frontend-asset/frontend.tar.gz /tmp/
COPY /config/config.yaml /app/config/config.yaml

# 创建数据目录
RUN mkdir -p /app/data

# 解压对应平台的 ech0 二进制
RUN mkdir -p /app/template && \
    if [ "$TARGETPLATFORM" = "linux/amd64" ]; then \
       tar -xzf /tmp/ech0-linux-amd64.tar.gz -C /app; \
    elif [ "$TARGETPLATFORM" = "linux/arm64" ]; then \
       tar -xzf /tmp/ech0-linux-arm64.tar.gz -C /app; \
    elif [ "$TARGETPLATFORM" = "linux/arm/v7" ]; then \
       tar -xzf /tmp/ech0-linux-armv7.tar.gz -C /app; \
    else \
       echo "Unsupported platform: $TARGETPLATFORM" && exit 1; \
    fi && \
    # 解压前端静态资源到 /app/template
    tar -xzf /tmp/frontend.tar.gz -C /app/template && \
    # 清理临时文件
    rm -rf /tmp/*.tar.gz

EXPOSE 6277

CMD ["/app/ech0"]
