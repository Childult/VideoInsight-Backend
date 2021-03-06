from audio_analysis.xunfei.lfasr import RequestApi


def audio_to_text(file: str) -> str:
    """
    语音转文本
    :param file: mp3文件路径
    :return: json字符串
    """
    xf_api = RequestApi(upload_file_path=file)
    return xf_api.all_api_request()
