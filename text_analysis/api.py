from audio_analysis.api import audio_to_text
from text_analysis.preprocess import sen_generator, preprocess_audio_text
from text_analysis.tencent import auto_summarization


def generate_abstract_from_audio(file: str) -> str:
    """
    根据音频文件，先转换为文本，然后进行文本摘要
    :param file:
    :return:
    """
    ret = {
        'AText': '',  # 语音转文本结果
        'TAbstract': '',  # 文本摘要
        'Error': ''  # 错误信息
    }
    text = preprocess_audio_text(audio_to_text(file))
    if text != '':
        ret['AText'] = text
    else:
        ret['Error'] = 'failed to transfer audio to text'

    ret['TAbstract'] = text_summarize(text)
    return str(ret)


def text_summarize(text: str) -> str:
    """
    文本摘要
    :param text: 经过预处理后的文本
    :return: 摘要
    """
    try:
        ret = []
        for seg in sen_generator(text):
            # 每次取出长度不超过2000字符的段落进行摘要
            ret.append(auto_summarization(seg))
            return ''.join(ret)
    except Exception as e:
        return str(e)


def cal_text_similarity(text_a, text_b: str) -> float:
    """
    计算两个文本相似度
    :param text_a: 文本1
    :param text_b: 文本2
    :return: 相似度[0, 1]，数值越大越相似
    """
    return 0.0


def text_search(document: str, keywords: [str]) -> [str]:
    """
    文本检索
    :param document: 文本
    :param keywords: 关键词列表
    :return: 文本中最相关的部分
    """
    return ['Not implemented.']
