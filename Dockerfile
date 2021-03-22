FROM python:3.7.10-slim

WORKDIR /app

EXPOSE 8080

COPY requirements.txt .

COPY bin/video_insight .

RUN chmod 777 video_insight \
    && apt-get update \
    && apt-get install -y --no-install-recommends ffmpeg \
    && ffmpeg -version \
    && rm -rf /var/lib/apt/lists/* \
    && mkdir -p /swc/code /swc/log /swc/resource/compressed \
    && pip install --no-cache-dir -r requirements.txt

CMD ["./video_insight"]
