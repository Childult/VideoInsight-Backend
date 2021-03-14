# 编译构建阶段
FROM golang:1.15.10-alpine as builder

COPY api_server /build/

WORKDIR /build

RUN go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w GO111MODULE=on \
    && GOOS=linux GOARCH=amd6 go build -v -ldflags "-X main.version=1.0" -o api_server

# 运行阶段
FROM python:3.7.10-slim

# 从编译阶段的中拷贝编译结果到当前镜像中
COPY --from=builder /build/api_server /swc/code

WORKDIR /swc/code

EXPOSE 8080

COPY requirements.txt /swc/code

RUN mkdir -p /swc/log /swc/resource/compressed \
    && pip install --no-cache-dir -r requirements.txt

ENTRYPOINT ["/api_server"]
