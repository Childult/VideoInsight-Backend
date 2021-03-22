# 编译构建阶段
FROM golang:1.16-alpine as builder

COPY api_server /build/

WORKDIR /build

RUN go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w GO111MODULE=on \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags "-X main.version=1.0" -o video_insight main.go

# 运行阶段
FROM alpine:latest

EXPOSE 8080

RUN apk --no-cache add ca-certificates \
    && mkdir -p /swc/code /swc/log /swc/resource/compressed

WORKDIR /swc/go

# 从编译阶段的中拷贝编译结果到当前镜像中
COPY --from=builder /build/video_insight .

CMD ["./video_insight"]
