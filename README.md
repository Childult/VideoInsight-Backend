# 提要钩玄APP后端

## 目录结构

- `api_server/`: RESTful API server
- `audio_analysis/`: 音频分析模块
- `text_analysis/`: 文本分析模块
- `video_analysis/`: 视频分析模块
- `video_getter/`: 视频获取模块

## 使用Docker运行

```bash
# 启动容器
docker-compose -p VideoInsight up -d

# 停止容器
docker-compose -p VideoInsight down
```