# 提要钩玄APP后端

## 目录结构

- `api_server/`: RESTful API server
- `audio_analysis/`: 音频分析模块
- `text_analysis/`: 文本分析模块
- `video_analysis/`: 视频分析模块
- `video_getter/`: 视频获取模块

## 环境需求

```bash
conda create -n SWC python=3.7
conda activate SWC
pip install -r requirements.txt
```

## 使用Docker运行

```bash
# 构建镜像
docker build -t swc .

# 启动容器
bash start-up.sh
```