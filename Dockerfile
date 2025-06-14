# =================== 构建阶段 ===================
FROM alpine:latest as builder

WORKDIR /app
ENV TZ=Asia/Shanghai

# 获取构建平台信息
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

# 创建数据目录
RUN mkdir -p /app/data && \
    mkdir -p /app/template && \
    mkdir -p /app/config

# 将所有平台的 ech0 二进制和前端资源复制进镜像
COPY /backend-artifacts/* /tmp/
COPY /frontend-asset/frontend.tar.gz /tmp/
COPY /config/config.yaml /app/config/config.yaml

# 解压对应平台的 ech0 二进制
RUN mkdir -p /app/template && \
    if [ "$TARGETOS" = "linux" ] && [ "$TARGETARCH" = "amd64" ]; then \
       tar -xzf /tmp/ech0-linux-amd64.tar.gz -C /app; \
    elif [ "$TARGETOS" = "linux" ] && [ "$TARGETARCH" = "arm64" ]; then \
       tar -xzf /tmp/ech0-linux-arm64.tar.gz -C /app; \
    elif [ "$TARGETOS" = "linux" ] && [ "$TARGETARCH" = "arm" ] && [ "$TARGETVARIANT" = "v7" ]; then \
       tar -xzf /tmp/ech0-linux-armv7.tar.gz -C /app; \
    else \
       echo "Unsupported platform: $TARGETOS/$TARGETARCH$TARGETVARIANT" && exit 1; \
    fi && \
    # 解压前端静态资源到 /app/template
    tar -xzf /tmp/frontend.tar.gz -C /app/template && \
    # 清理临时文件
    rm -rf /tmp/*

# =================== 最终镜像 ===================
FROM alpine:latest

WORKDIR /app
ENV TZ=Asia/Shanghai

COPY --from=builder /app /app

EXPOSE 6277

CMD ["/app/ech0"]
