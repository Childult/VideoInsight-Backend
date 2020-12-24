# 提要钩玄APP后端

## 目录结构

- `swagger_server/`: flask server，由[swagger-codegen](https://github.com/swagger-api/swagger-codegen)生成
- `audio_analysis/`: 音频分析模块
- `text_analysis/`: 文本分析模块
- `video_analysis/`: 视频分析模块
- `video_getter/`: 视频获取模块

## 环境需求

```bash
conda create -n SWC python=3.6
conda activate SWC
pip install -r requirements.txt
```

## 启动flask server

```bash
conda activate SWC
python -m swagger_server
```

打开浏览器访问：[http://localhost:9090/v2/ui/](http://localhost:9090/v2/ui/)

## 使用Docker运行

```bash
# 构建镜像
docker build -t swagger_server .

# 启动容器
docker run -p 9090:9090 swagger_server
```