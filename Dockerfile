# 使用 Node.js 镜像作为前端构建阶段
FROM node:22.13.0 AS frontend-build

# 设置前端工作目录
WORKDIR /app/web

# 复制前端代码
COPY ./web/package.json ./web/package-lock.json* ./
COPY ./web/ .
RUN npm install  # 安装前端依赖

# 构建前端项目
RUN npm run generate  # 生成静态文件

# 复制构建后的文件到后端 public 目录
RUN cp -r .output/public /app/public/

# 使用 Golang 镜像作为后端构建阶段
FROM golang:1.24-alpine AS backend-build

# 设置后端工作目录
WORKDIR /app

# 安装必要的依赖，尤其是编译 cgo 所需的工具（如 gcc）
RUN apk add --no-cache gcc musl-dev

# 设置环境变量启用 cgo
ENV CGO_ENABLED=1

# 复制后端代码和必要文件
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg
COPY ./config ./config
# COPY ./public ./public

# 创建目录
RUN mkdir -p /app/data && chmod -R 777 /app/data

# 构建 Go 后端应用
RUN go build -o /app/ech0 ./cmd/server/main.go

# 使用 Golang 镜像作为运行时阶段
FROM golang:1.24-alpine

WORKDIR /app

# 复制构建后的二进制文件和前端资源
COPY --from=backend-build /app/config /app/config
COPY --from=backend-build /app/ech0 /app/ech0
COPY --from=frontend-build /app/public /app/public

# 暴露端口
EXPOSE 1314

# 运行后端服务
CMD ["/app/ech0"]
