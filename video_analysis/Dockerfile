FROM python:3.7
LABEL maintainer="Zhu Jie<whuzhujie@qq.com>"

RUN pip install -U pip
RUN pip config set global.index-url https://mirror.baidu.com/pypi/simple
RUN pip config set install.trusted-host mirror.baidu.com
RUN pip --default-timeout=1000 install -U pip

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

COPY requirements.txt /usr/src/app/

RUN pip install --no-cache-dir -r requirements.txt
RUN pip install torch==1.7.1+cpu torchvision==0.8.2+cpu torchaudio===0.7.2 -f https://download.pytorch.org/whl/torch_stable.html
RUN pip install opencv-python==4.4.0.46
RUN pip install h5py==3.1.0

COPY . /usr/src/app

EXPOSE 9091

ENTRYPOINT ["python"]

# CMD ["-m", "video_analysis"]