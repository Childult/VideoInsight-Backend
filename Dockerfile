FROM python:3.7

RUN apt-get update \
    && apt-get install -y sudo \
    && sudo apt update \
    && sudo apt install -y libgl1-mesa-glx \
    && sudo apt install -y ffmpeg \
    && rm -rf /var/lib/apt/lists/* \
    && mkdir -p /swc/code /swc/log /swc/resource /swc/resource/compressed \
    && wget https://golang.google.cn/dl/go1.15.8.linux-amd64.tar.gz \
    && rm -rf /usr/local/go \
    && tar -C /usr/local -xzf go1.15.8.linux-amd64.tar.gz \
    && rm -f go1.15.8.linux-amd64.tar.gz \
    && export PATH=$PATH:/usr/local/go/bin \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w GO111MODULE=on \
    && pip install torch==1.7.1+cpu torchvision==0.8.2+cpu torchaudio===0.7.2 -f https://download.pytorch.org/whl/torch_stable.html \
    && pip install opencv-python==4.4.0.46 \
    && pip install h5py==3.1.0 \
    && wget -P /root/.cache/torch/hub/checkpoints/ https://download.pytorch.org/models/resnet152-b121ed2d.pth

EXPOSE 8080

WORKDIR /swc/code

COPY requirements.txt /swc/code

RUN pip install --no-cache-dir -r requirements.txt

CMD ["bash"]
