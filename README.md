# Ech0s

developing

## docker-componse部署
创建一个新目录并将 docker-compose.yml 文件放入其中
在该目录下执行以下命令启动服务：
```shell
docker-compose up -d
```

## docker部署

```shell
docker run -d \
  --name ech0 \
  -p 1314:1314 \
  -v /opt/ech0/data:/app/data \
  sn0wl1n/ech0:v2.2.0
```
