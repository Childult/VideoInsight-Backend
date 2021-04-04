#!/usr/bin/env sh
set -x

docker image rm videoinsight_api-server:latest

# 启动容器
docker-compose -p VideoInsight up -d

echo "Start the backend successfully!"
