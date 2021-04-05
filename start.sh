#!/usr/bin/env sh
set -x

rm -f ../swc-log/* > /dev/null

docker image rm videoinsight_api-server:latest

chmod +x bin/video_insight

# 启动容器
docker-compose -p VideoInsight up -d

echo "Start the backend successfully!"
