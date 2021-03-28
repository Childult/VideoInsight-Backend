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
        print(seg)
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
    print(text_summarize(
    '这个自然语言处理面对的面临的三大真实挑战，实际上是就现实挑战，实际挑战倒不这个这个词用的不一定。特别对，反正就现在我们面临的三大真实挑战，这个自然源处理的话被称作这个人工智能皇冠上的明珠。这个各位这个大佬哈都有很多论述，比如说这个图灵奖得主哈燕乐村就说深度学习的下一个前沿课题是自然语言理解等等，吧有很多这种论述。然后我个人也认为以自然语言为核心的语义理解哈是机器难以逾越的鸿沟。因为这个语言这一关如果机器搞明白了以后，这个我的话就机器真的成了精了，那时候我们人类真有危险了，真有危险，现在还还还还不是，所以这就进进一步彰显这个自然语言处理的这个困难，它已经成为制约人工智能取得更大突破的主要瓶颈之一，这是我的一个基本判断。那自然源处理的话，它实际上历史上有两大研究范式，就是理性主义、经验主义、经验主义我这就不讲了。实际上最近这最近这个经验主义的话，从90年代一直到现在也分了几波，最近这一两从18年到现在，实际上是我们大家都都熟悉的，就大规模预训练语言模型，就是GD pGDP什么board，反正这一路到GDP three这一套东西，那这个 Gptc的话，今天其实大家很多都在谈，可见这个事对我们的这个震动，我的理解它实际上就是几个大就极大，即大规模模型及大规模数据集大规模计算，三个吉大光大还不够，加个急就导致出这个效果。这个给人感觉好像有量变引起质变的这个感觉，这个一方面是它的性能超乎我们的想象，一般我们想好像规模大不会引起这种变化。第二他也有一些这种一些科学现象，觉得有些奇怪，比如说这个这个这个这个这个就就是那个叫什么？Double descent。Uh那个 Dwd3的现象也很奇怪，这个跟我们一般机器学习那个道理不太一样，反正也所以它有量变引起质变的这么一些趋势。这个你去看，我这个我一直举一个例子，你去看关羽，哈这个 Weekday里边关羽，它只有一些最简单的关于的这个这个这个属性描写，关于是是一个人，关羽是个将军，关于4个蜀国的这个关于他儿子是谁，孩子是谁是吧？就所有关于关羽的什么那个什么过五关斩六将，什么这个这个三英战吕布等等的所有事通通没有统统没有，所以你这种这种知识库大，它其实它它覆盖面很窄。是吧？'))
    print(time.time() - beg)
