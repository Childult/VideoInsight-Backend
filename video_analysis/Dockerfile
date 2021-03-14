FROM pytorch/pytorch:1.7.0-cuda11.0-cudnn8-runtime

EXPOSE 50051

RUN apt-get update \
    && apt-get install -y sudo wget libglib2.0-0 \
    && sudo apt update \
    && sudo apt install -y libgl1-mesa-glx \
    && rm -rf /var/lib/apt/lists/* \
    && pip install opencv-python==4.4.0.46 \
    && python -c "import cv2 ; print(cv2.__version__)" \
    && mkdir -p /swc/code/video_analysis /swc/log /swc/resource/compressed \
    && pip install h5py==3.1.0 \
    && wget -P /root/.cache/torch/hub/checkpoints/ https://download.pytorch.org/models/resnet152-b121ed2d.pth

WORKDIR /swc/code/video_analysis

COPY requirements.txt /swc/code/video_analysis

RUN pip install --no-cache-dir -r requirements.txt

CMD ["python", "main.py"]
