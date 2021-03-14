import os

from audio_analysis.xunfei.lfasr import RequestApi


def audio_to_text(file: str) -> dict:
    """
    语音转文本
    :param file: mp3文件路径
    :return: 包含结果的键值对，key=['data', 'ok']
    """
    xf_api = RequestApi(upload_file_path=file)
    return xf_api.all_api_request()


def extract_audio(file: str) -> str:
    """
    提取视频中的音频文件
    :param file: 视频文件
    :return: mp3文件名或者错误信息
    """
    title = os.path.basename(file).split('.', -1)[0]
    title_path = file.split('.', -1)[0]
    try:
        os.system('ffmpeg -i ' + file + ' -f mp3 -ar 16000 ' + title_path + '.mp3')
    except Exception as e:
        return "Error: {0}".format(str(e))
    return title + '.mp3'
