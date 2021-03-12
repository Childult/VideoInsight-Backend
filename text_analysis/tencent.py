# -*- encoding:utf-8 -*-
import json
from tencentcloud.common import credential
from tencentcloud.common.profile.client_profile import ClientProfile
from tencentcloud.common.profile.http_profile import HttpProfile
from tencentcloud.common.exception.tencent_cloud_sdk_exception import TencentCloudSDKException
from tencentcloud.nlp.v20190408 import nlp_client, models


def auto_summarization(text: str, length_ratio: float = 0.3) -> str:
    """
    利用人工智能算法，自动抽取文本中的关键信息并生成指定长度的文本摘要。
    :param text: 待摘要的文本
    :param length_ratio: 生成摘要的长度相比于原文长度的占比
    :return: 生成的摘要
    """
    try:
        cred = credential.Credential("AKIDCF5h4DJrLcgW90p9UgeOvFLEsIUUao9G", "1oSt54pNhnfPwFKBpFRQBj4D9rMdhX9X")
        httpProfile = HttpProfile()
        httpProfile.endpoint = "nlp.tencentcloudapi.com"

        clientProfile = ClientProfile()
        clientProfile.httpProfile = httpProfile
        client = nlp_client.NlpClient(cred, "ap-guangzhou", clientProfile)

        req = models.AutoSummarizationRequest()
        params = {
            "Length": int(length_ratio * len(text)),
            "Text": text
        }
        req.from_json_string(json.dumps(params))

        resp = client.AutoSummarization(req)
        return json.loads(resp.to_json_string())['Summary']

    except TencentCloudSDKException as err:
        raise err


if __name__ == '__main__':
    doc = """
    这些事件发生时，我刚从美国内布拉斯加州的贫瘠地区做完一项科考工作回来。我当时是巴黎自然史博物馆的客座教授，法国政府派我参加这次考察活动。我在内布拉斯加州度过了半年时间，收集了许多珍贵资料，满载而归，3 月底抵达纽约。我决定 5 月初动身回法国。于是，我就抓紧这段候船逗留时间，把收集到的矿物和动植物标本进行分类整理，可就在这时，斯科舍号出事了。
我对当时的街谈巷议自然了如指掌，再说了，我怎能听而不闻、无动于衷呢？我把美国和欧洲的各种报刊读了又读，但未能深入了解真相。神秘莫测，百思不得其解。我左思右想，摇摆于两个极端之间，始终形不成一种见解。其中肯定有名堂，这是不容置疑的，如果有人表示怀疑，就请他们去摸一摸斯科舍号的伤口好了。
我到纽约时，这个问题正炒得沸反盈天。某些不学无术之徒提出设想，有说是浮动的小岛，也有说是不可捉摸的暗礁，不过，这些个假设通通都被推翻了。很显然，除非这暗礁腹部装有机器，不然的话，它怎能如此快速地转移呢？
同样的道理，说它是一块浮动的船体或是一堆大船残片，这种假设也不能成立，理由仍然是移动速度太快。
那么，问题只能有两种解释，人们各持己见，自然就分成观点截然不同的两派：一派说这是一个力大无比的怪物，另一派说这是一艘动力极强的“潜水船”。
哦，最后那种假设固然可以接受，但到欧美各国调查之后，也就难以自圆其说了。有哪个普通人会拥有如此强大动力的机械？这是不可能的。他在何地何时叫何人制造了这么个庞然大物，而且如何能在建造中做到风声不走漏呢？
看来，只有政府才有可能拥有这种破坏性的机器，在这个灾难深重的时代，人们千方百计要增强战争武器威力，那就有这种可能，一个国家瞒着其他国家在试制这类骇人听闻的武器。继夏斯勃步枪之后有水雷，水雷之后有水下撞锤，然后魔道攀升反应，事态愈演愈烈。至少，我是这样想的。
    """
    print(auto_summarization(doc.strip(), 200))
