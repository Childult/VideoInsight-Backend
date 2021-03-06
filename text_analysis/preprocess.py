import json

from textrank4zh import TextRank4Sentence

stopwords = [
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
]


def preprocess_audio_text(json_str: str) -> str:
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


def text_rank_summarize(text: str, num: int = 3) -> []:
    """
    基于text rank的文本摘要。
    :param text: 要生成摘要的文本
    :param num: 摘要包含的句子数目
    :return:
    """
    tr4s = TextRank4Sentence()
    tr4s.analyze(text=text, lower=True, source='all_filters')
    return tr4s.get_key_sentences(num=num)
    # print(item.index, item.weight, item.sentence)  # index是语句在文本中位置，weight是权重


def extract_document(filename: str) -> str:
    """
    从讯飞API返回的json文本中提取语音识别的结果，并删除其中的停用词。
    :param filename: json文件名
    :return: 删除停用词后的字符串
    """
    ret = []
    with open(filename, 'r') as f:
        data = f.read()
        js = json.loads(data)
        for entry in js:
            text = entry.get('onebest')
            if text:
                ret.append(text.strip())
    return remove_stopwords(''.join(ret))


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


def replace_full_width_punc(s: str) -> str:
    """
    将字符串中的部分全角标点替换为半角标点。
    不替换全角句号，因为半角句号可能是小数点。
    :param s:
    :return:
    """
    full = u'，！？【】（）％＃＠＆'
    half = u',!?[]()%#@&'
    trantab = str.maketrans(half, full)
    return s.translate(trantab)


def sentence_split(text: str):
    """
    中文分句，判断依据为是否以句号、问号或感叹号结尾。
    :param text:
    :return:
    """
    text = replace_full_width_punc(text)
    ret = []
    for sen in text.split('。'):
        if '?' in sen:
            ret.extend(sen.split('?'))
        elif '!' in sen:
            ret.extend(sen.split('!'))
        else:
            ret.append(sen)
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


if __name__ == '__main__':
    t = extract_document('AText.json')
    for seg in sen_generator(t):
        print(len(seg), seg)
    # text = None
    # with open('text.txt', 'r') as f:
    #     text = f.read()

    # old_time = time.time()
    # l2 = text_rank(t)
    # print(time.time() - old_time)
    #
    #
    # old_time = time.time()
    # po = FastTextRank4Sentence(use_w2v=False, tol=0.0001).summarize(t, 3)
    # print(time.time() - old_time)
    #
    # for s in lst:
    #     print(s)
    # for s in l2:
    #     print(s)
    # for s in po:
    #     print(s)
