import sys
import os
import you_get
import json
import time
import base64
import signal, functools


def get_video_title(url: str) -> str:
    """
    获取视频标题
    :param url: 视频URL
    :return: 若成功返回视频标题，否则返回错误信息
    """
    try:
        cmd = 'you-get ' + url + ' --json'
        r = os.popen(cmd)
        text = r.buffer.read().decode(encoding='utf8')
        parsed_text = json.loads(text)
        title = parsed_text['title']
        res = title
    except Exception as e:
        res = e
    return res


def download_video(url: str, path: str, timeout: int=7200) -> str:
    """
    下载视频
    :param url: 视频URL
    :param path: 视频存储路径
    :param timeout: 下载超时时间
    :return: 若成功返回文件名，否则返回错误信息
    """
    cmd = 'you-get ' + url + ' --json'
    r = os.popen(cmd)
    text = r.buffer.read().decode(encoding='utf8')
    parsed_text = json.loads(text)
    
    # get title
    ts = time.time()
    title = base64.urlsafe_b64encode((str(ts) + url).encode('utf8')).decode('utf-8')

    # get accessible qualities (480p default), format and container
    qualities = ['480', '720', '360', '240', '1080']
    for quality in qualities:
        formats = [k for k, v in parsed_text['streams'].items() if quality in v['quality']]
        if formats:
            break
    if not formats:
        formats = [k for k, v in parsed_text['streams'].items()]
    formats.sort()
    fmt = formats[0]
    container = parsed_text['streams'][fmt]['container']

    sys.argv = ['you-get', '-o', path, '-O', title, url, '--format=%s' % fmt, '--debug']
    file = title + '.' + container

    def handler(signum, frame):
        raise TimeoutError()
    
    signal.signal(signal.SIGALRM, handler)
    signal.alarm(timeout)

    try:
        you_get.main()
    except Exception as e:
        return "Error: {0}".format(e)
    finally:
        signal.alarm(0)
        signal.signal(signal.SIGALRM, signal.SIG_DFL)
        if not container == 'mp4':
            os.system('ffmpeg -i ' + path + '/' + file + ' -vf scale=480:-1' + ' -vcodec h264 ' + title + '.mp4')
            file = title + '.mp4'
    
    return file


if __name__ == '__main__':
    test_url = '<test url>'
    download_video(test_url, '.')
    print(get_video_title(test_url))
