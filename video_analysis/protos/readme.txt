# 与API Server之间的RPC协议
python -m grpc_tools.protoc -I video_analysis/protos --python_out=video_analysis/ --grpc_python_out=video_analysis/ video_analysis/protos/video.proto
