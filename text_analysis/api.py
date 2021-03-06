import json

from text_analysis.preprocess import remove_stopwords, sen_generator
from text_analysis.tencent import auto_summarization


def preprocess(json_str: str) -> str:
    """
    对讯飞语音转文本的结果进行预处理
    :param json_str: json字符串
    :return: 返回删除停用词的字符串
    """
    js = json.loads(json_str)
    if js['ok'] == 0:
        ret = []
        for entry in js['data']:
            text = entry.get('onebest')
            if text:
                ret.append(text.strip())
        return remove_stopwords(''.join(ret))
    else:
        return ''


def text_summarize(text: str) -> str:
    """
    文本摘要
    :param text: 经过预处理后的文本
    :return: 摘要
    """
    try:
        ret = []
        for seg in sen_generator(text):
            print(len(seg), seg)
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
