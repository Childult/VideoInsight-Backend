# -*- encoding:utf-8 -*-
import json
from tencentcloud.common import credential
from tencentcloud.common.profile.client_profile import ClientProfile
from tencentcloud.common.profile.http_profile import HttpProfile
from tencentcloud.common.exception.tencent_cloud_sdk_exception import TencentCloudSDKException
from tencentcloud.nlp.v20190408 import nlp_client, models


class TencentCloudNLP:
    secrete_id = 'AKIDCF5h4DJrLcgW90p9UgeOvFLEsIUUao9G'
    secrete_key = '1oSt54pNhnfPwFKBpFRQBj4D9rMdhX9X'
    endpoint = 'nlp.tencentcloudapi.com'
    region = 'ap-guangzhou'

    def __init__(self):
        cred = credential.Credential(TencentCloudNLP.secrete_id, TencentCloudNLP.secrete_key)
        http_profile = HttpProfile()
        http_profile.endpoint = TencentCloudNLP.endpoint

        client_profile = ClientProfile()
        client_profile.httpProfile = http_profile
        self.client = nlp_client.NlpClient(cred, TencentCloudNLP.region, client_profile)

    def auto_summarization(self, text: str, length_ratio: float = 0.2) -> str:
        """
        利用人工智能算法，自动抽取文本中的关键信息并生成指定长度的文本摘要。
        :param text: 待摘要的文本
        :param length_ratio: 生成摘要的长度相比于原文长度的占比
        :return: 生成的摘要
        """
        try:

            req = models.AutoSummarizationRequest()
            params = {
                "Length": int(length_ratio * len(text)),
                "Text": text
            }
            req.from_json_string(json.dumps(params))

            resp = self.client.AutoSummarization(req)
            return json.loads(resp.to_json_string())['Summary']

        except TencentCloudSDKException as err:
            raise err

    def text_correction(self, text: str) -> str:
        """
        文本纠错
        :param text:
        :return:
        """
        try:
            req = models.TextCorrectionRequest()
            params = {
                "Text": text
            }
            req.from_json_string(json.dumps(params))

            resp = self.client.TextCorrection(req)
            json_obj = json.loads(resp.to_json_string())
            # print(json_obj['CCITokens'])
            return json_obj['ResultText']
        except TencentCloudSDKException as err:
            raise err

    def sentence_vector(self, text: str) -> str:
        """
        句子向量
        :param text:
        :return:
        """
        try:
            req = models.SentenceEmbeddingRequest()
            params = {
                "Text": text
            }
            req.from_json_string(json.dumps(params))

            resp = self.client.SentenceEmbedding(req)
            return json.loads(resp.to_json_string())['Vector']
        except TencentCloudSDKException as err:
            raise err

    def text_similarity(self, src: str, tar: str) -> float:
        """
        文本相似度
        :param src:
        :param tar:
        :return:
        """
        try:
            req = models.TextSimilarityRequest()
            params = {
                "SrcText": src,
                "TargetText": [tar]
            }
            req.from_json_string(json.dumps(params))

            resp = self.client.TextSimilarity(req)
            return json.loads(resp.to_json_string())['Similarity'][0]['Score']
        except TencentCloudSDKException as err:
            raise err


TENCENT_API = TencentCloudNLP()
