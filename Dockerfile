FROM swc_python:latest

WORKDIR /app

EXPOSE 8080

COPY bin/video_insight .

CMD ["./video_insight"]
