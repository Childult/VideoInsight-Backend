FROM python:3.7
LABEL maintainer="Zhu Jie<whuzhujie@qq.com>"

RUN apt-get update \
    && apt-get install -y sudo \
    && sudo apt update \
    && sudo apt install -y libgl1-mesa-glx \
    && rm -rf /var/lib/apt/lists/* \
    && mkdir -p /usr/src/app

WORKDIR /usr/src/app

COPY requirements.txt /usr/src/app/

RUN pip install --no-cache-dir -r requirements.txt \
    && pip install torch==1.7.1+cpu torchvision==0.8.2+cpu torchaudio===0.7.2 -f https://download.pytorch.org/whl/torch_stable.html \
    && pip install opencv-python==4.4.0.46 \
    && pip install h5py==3.1.0 \
    && wget -P /root/.cache/torch/hub/checkpoints/ https://download.pytorch.org/models/resnet152-b121ed2d.pth

CMD ["bash"]
