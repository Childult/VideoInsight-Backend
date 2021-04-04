# -*- encoding:utf-8 -*-
import json
import re

from text_analysis.tencent import TENCENT_API

stopwords = [
    '\n',
    '嗯',
    "呃",
    "啊",
    "呗",
    "咦",
    "喏",
    "喔",
    "唷",
    "嗬",
    "嗳",
    "呢",
    "哈",
    "嘛",
    "哎",
    "唉",
    "呀",
    "吧",
]

api = TENCENT_API


def preprocess_audio_text(result: dict) -> str:
    """
    对讯飞语音转文本的结果进行预处理
    :param result: 识别结果键值对
    :return: 返回删除停用词的字符串
    """
    if result['ok'] == 0:
        js = json.loads(result['data'])
        ret = []
        for entry in js:
            text = entry.get('onebest')
            if text:
                ret.append(text.strip())
        return remove_stopwords(''.join(ret))
    else:
        return ''


def remove_stopwords(t: str) -> str:
    """
    删除字符串中的停用词。
    :param t:
    :return:
    """
    ret = t
    for w in list(set(stopwords)):
        ret = ret.replace(w, '')
    return ret


def full_to_half(s: str) -> str:
    """
    将字符串中的部分全角标点替换为半角标点。
    :param s:
    :return:
    """
    ret = []
    for char in s:
        num = ord(char)
        if num == 0xFF0e:
            # 不替换全角句号为半角，因为半角句号可能是小数点。
            pass
        elif num == 0x3000:
            num = 32
        elif 0xFF01 <= num <= 0xFF5E:
            num -= 0xfee0
        ret.append(chr(num))
    return ''.join(ret).strip()


def sentence_split(text: str) -> list:
    """
    中文分句，判断依据为是否以句号、问号或感叹号结尾。
    :param text:
    :return:
    """
    ret = []
    count = -1
    for sen in re.split(r'([。?!])', full_to_half(text)):
        if len(sen) > 1:
            ret.append(sen)
            count += 1
        elif len(sen) == 1:
            ret[count] += sen
    return ret


def segment_split(sens: [str]) -> [str]:
    """
    计算两两句子的相似度，通过相似度的变化进行分段
    :param sens:
    :return:
    """
    # 数量太少，无法判断
    if len(sens) <= 2:
        return sens

    sim = []
    for i in range(len(sens) - 1):
        sim.append(api.text_similarity(sens[i], sens[i + 1]))

    ret = []
    precision = 0.01
    beg = 0
    for i in range(1, len(sim)):
        if sim[i - 1] - sim[i] > precision:
            ret.append(''.join(sens[beg:i + 1]))
            beg = i + 1

    if beg < len(sim):
        ret.append(''.join(sens[beg:]))

    return ret


def sen_generator(text: str, max_len=2000):
    """
    限定长度文本生成器，每次返回不超过最大长度的文本段落。
    :param text:
    :param max_len: 每次返回字符串的最大长度
    :return:
    """
    temp = []
    count = 0
    for sen in sentence_split(text):
        if count + len(sen) < max_len:
            temp.append(sen)
            count += len(sen)
        else:
            yield ''.join(temp)
            temp.clear()
            count = len(sen)
            temp.append(sen)
    yield ''.join(temp)
