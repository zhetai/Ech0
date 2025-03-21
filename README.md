# Ech0s

developing

## docker部署

```shell
docker run -d \
  --name ech0 \
  -p 1314:1314 \
  -v /opt/ech0/data:/app/data \
  sn0wl1n/ech0:v2.2.0
```
