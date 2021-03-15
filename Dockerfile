# 编译构建阶段
FROM golang:1.16-alpine as builder

COPY api_server /build/

WORKDIR /build

RUN go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w GO111MODULE=on \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags "-X main.version=1.0" -o video_insight main.go

# 运行阶段
FROM python:3.7.10-slim

WORKDIR /app

# 从编译阶段的中拷贝编译结果到当前镜像中
COPY --from=builder /build/video_insight .

EXPOSE 8080

COPY requirements.txt .

RUN apt-get update \
    && apt-get install -y --no-install-recommends ffmpeg \
    && ffmpeg -version \
    && rm -rf /var/lib/apt/lists/* \
    && mkdir -p /swc/code /swc/log /swc/resource/compressed \
    && pip install --no-cache-dir -r requirements.txt

CMD ["./video_insight"]
