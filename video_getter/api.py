import sys
import os
import you_get
import json
import time
import base64


def download_video(url: str, path: str) -> str:
    """
    下载视频
    :param url: 视频URL
    :param path: 视频存储路径
    :return: 若成功返回文件名，否则返回错误信息
    """
    cmd = 'you-get -o ' + path + ' ' + url + ' --json'
    r = os.popen(cmd)
    text = r.buffer.read().decode(encoding='utf8')
    parsed_text = json.loads(text)

    # get title
    ts = time.time()
    title = base64.b64encode((str(ts) + url).encode('utf8')).decode('utf-8')

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
    try:
        you_get.main()
    except Exception as e:
        return "Error: {0}".format(str(e))

    return title + '.' + container
