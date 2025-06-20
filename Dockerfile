# # =================== 构建阶段 ===================
# FROM alpine:latest AS builder

# WORKDIR /app
# ENV TZ=Asia/Shanghai

# # 获取构建平台信息
# ARG TARGETOS
# ARG TARGETARCH
# ARG TARGETVARIANT

# # 创建数据目录
# RUN mkdir -p /app/data && \
#     mkdir -p /app/template && \
#     mkdir -p /app/config

# # 将所有平台的 ech0 二进制和前端资源复制进镜像
# COPY /backend-artifacts/* /tmp/
# COPY /frontend-asset/frontend.tar.gz /tmp/
# COPY /config/config.yaml /app/config/config.yaml

# # 解压对应平台的 ech0 二进制
# RUN mkdir -p /app/template && \
#     if [ "$TARGETOS" = "linux" ] && [ "$TARGETARCH" = "amd64" ]; then \
#        tar -xzf /tmp/ech0-linux-amd64.tar.gz -C /tmp && mv /tmp/ech0-linux-amd64 /app/ech0; \
#     elif [ "$TARGETOS" = "linux" ] && [ "$TARGETARCH" = "arm64" ]; then \
#        tar -xzf /tmp/ech0-linux-arm64.tar.gz -C /tmp && mv /tmp/ech0-linux-arm64 /app/ech0; \
#     elif [ "$TARGETOS" = "linux" ] && [ "$TARGETARCH" = "arm" ] && [ "$TARGETVARIANT" = "v7" ]; then \
#        tar -xzf /tmp/ech0-linux-armv7.tar.gz -C /tmp && mv /tmp/ech0-linux-armv7 /app/ech0; \
#     else \
#        echo "Unsupported platform: $TARGETOS/$TARGETARCH$TARGETVARIANT" && exit 1; \
#     fi && \
#     # 解压前端静态资源到 /app/template
#     tar -xzf /tmp/frontend.tar.gz -C /app/template && \
#     # 清理临时文件
#     rm -rf /tmp/*

# # =================== 最终镜像 ===================
# FROM debian:bookworm-slim

# WORKDIR /app
# ENV TZ=Asia/Shanghai

# COPY --from=builder /app /app

# RUN ls -lh /app

# # 设置 ech0 二进制文件的权限
# RUN chmod +x /app/ech0

# EXPOSE 6277

# CMD ["/app/ech0"]

# =================== 构建阶段 ===================
# 使用与最终镜像一致的 alpine 作为构建器，保持环境统一
FROM alpine:latest AS builder

# 接收 Docker Buildx 自动提供的构建参数
ARG TARGETARCH

WORKDIR /app

# --- 简化文件复制 ---
# Docker 的构建上下文就是代码仓库的根目录，直接使用相对路径即可
# 将所有需要的文件一次性复制到构建器中
COPY backend-artifacts/ech0-linux-${TARGETARCH}.tar.gz /tmp/backend.tar.gz
COPY frontend-asset/frontend.tar.gz /tmp/frontend.tar.gz
COPY config/config.yaml /app/config/config.yaml

# --- 简化资源解压 ---
# 利用 TARGETARCH 变量，无需 if/else 判断，直接解压对应的文件
RUN tar -xzf /tmp/backend.tar.gz -C /app && \
    mkdir -p /app/template && \
    tar -xzf /tmp/frontend.tar.gz -C /app/template && \
    # 重命名二进制文件为统一的名称 'ech0'
    mv /app/ech0-linux-${TARGETARCH} /app/ech0 && \
    # 增加执行权限
    chmod +x /app/ech0

# =================== 最终镜像 ===================
# 使用 alpine 作为最终镜像，体积小，启动快
FROM alpine:latest

WORKDIR /app

# 从构建器中只复制最终需要的文件，保持镜像干净
COPY --from=builder /app /app

# 再次确认权限（虽然上一步已经做过，但这是一个好习惯）
RUN chmod +x /app/ech0 && \
    # 添加一个验证步骤，确保文件确实存在且可执行
    ls -lh /app/ech0

EXPOSE 6277

# 运行程序
CMD ["/app/ech0"]