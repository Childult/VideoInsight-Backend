# -*- encoding:utf-8 -*-
import json
import time

from audio_analysis.api import audio_to_text
from text_analysis.utils import sen_generator, preprocess_audio_text, segment_split, sentence_split
from text_analysis.tencent import TENCENT_API


def generate_abstract_from_audio(file: str) -> str:
    """
    根据音频文件，先转换为文本，然后进行文本摘要
    :param file:
    :return:
    """
    ret = {
        'AText': '',  # 语音转文本结果
        'TAbstract': None,  # 文本摘要
        'Error': ''  # 错误信息
    }
    try:
        text = preprocess_audio_text(audio_to_text(file))
    except Exception as e:
        print(e)
        ret['Error'] = 'failed to transfer audio to text'
    else:
        try:
            if text != '':
                ret['AText'] = text
                ret['TAbstract'] = text_summarize(text)
        except Exception as e:
            print(e)
            ret['Error'] = 'failed to summarize the text'

    return json.dumps(ret, ensure_ascii=False)


def text_summarize(text: str) -> [str]:
    """
    文本摘要
    :param text: 经过预处理后的文本
    :return: 摘要列表，按照语义进行了切分，每个列表元素是一个段落
    """
    ret = []
    for seg in sen_generator(text):
        try:
            # 文本纠错
            tmp = TENCENT_API.text_correction(seg)
            # 每次取出长度不超过2000字符的段落进行摘要
            ret.append(TENCENT_API.auto_summarization(tmp))
        except Exception as e:
            print('in text_summarize: ', e)
    return segment_split(ret)


def cal_text_similarity(text_a, text_b: str) -> float:
    """
    计算两个文本相似度
    :param text_a: 文本1
    :param text_b: 文本2
    :return: 相似度[0, 1]，数值越大越相似
    """
    return TENCENT_API.text_similarity(text_a, text_b)


def text_search(document: str, keywords: [str]) -> [str]:
    """
    文本检索
    :param document: 文本
    :param keywords: 关键词列表
    :return: 文本中最相关的部分
    """
    return ['Not implemented.']


if __name__ == '__main__':
    beg = time.time()
    print(generate_abstract_from_audio('/swc/resource/1616397097/MTYxNjM5NzEzNy4yNDI3MjhodHRwczovL3d3dy5iaWxpYmlsaS5jb20vdmlkZW8vQlYxQnY0MTFrNzQ1.mp3'))
    print(time.time() - beg)
