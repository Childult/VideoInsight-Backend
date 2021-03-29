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

if __name__ == '__main__':
    # api = TencentCloudNLP()
    # doc = """第一个问题，就是这个表达能力高的表达能力到底是不是必要的？我们有有没有必要构造一个这么高血压能力的一个出神经网络？当我们知道这个表达能力，给我们理解这个继人的这个能力提供了一个理论上的就他就知道你什么时候你能干什么事不能干。比如说刚才讲的那个用那个分子结构来预测一个化学分子的那个它的量子特性这个事情，他是你他肯定是搞不定的。他在有些场景下可以有些简单的统计量可以统计，但是有些他搞不定，然后就这样的话你就知道哪些任务你也不用去试，是你也做不好，哪些任务你是可以可以通过试的方式来来解决的。"""
    # sens = sentence_split(doc)
    # lst = []
    sens = ['各位同学大家好,对于刑法学的唉呀今年3月1号是一个大日子,因为刑法修正案11生效了,刑法中增加了许多新的罪名,法网被进一步的严密。',
            '所以今天想给大家简单介绍一下,修正案11主要增加了哪些新罪名,哪些行为已经开始变成犯罪了?', '我先想问大家一个问题,大家觉得刑法修正案11是扩张的刑罚权还是收缩了刑罚权?',
            '扩张咱选一,收缩咱们选答案,当然是老规矩,应该选商刑法修正案,其实系扩张的刑罚权有收缩的刑罚权,所以这是有些新增的罪名属于真正的新增。', '也就是说以前这些行为不构成犯罪,3月1号之后才构成犯罪的。',
            '那根据从旧兼从轻原则,如果一种新的法律对行为人不利,那就不能溯及既往,毕竟刑法嘛是一根带哨子的鞭子,在抽你之前是要告诉你的,那么这些让无罪变为有罪的真正意义上的新尊罪名有哪些?',
            '那主要有如下一些,一个是商业间谍行为,为境外的机构组织人员窃取、刺探、收买非法提供商业秘密的,但是可以处5年以下有期徒刑,那么情节特别严重的,甚至还可以触到5年以上,最高可以判到15年有期徒刑。',
            '因为刑法以前只有为境外窃取刺探收买非法提供国家秘密,但现在商业间谍是越来越多,所以就增加了这个犯罪比茅台酒的配方加工的独特方法。', '虽然茅台酒厂是国企,但这个肯定不属于国家秘密,嘛但它属于商业秘密。',
            '那如果张三把这个独特的加工方法告诉了外国的酒厂,那就有可能构成。',
            '那第二是负有照护职责,人员性侵罪这个我们以前也提过,对已满14周岁不满16周岁的未成年女性负有监护、收养、看护、教育、医疗等特殊职责的人员,与该未成年女性发生性关系处三年以下有期徒刑,那情节恶劣,可以出3~10年。',
            '那大家知道以前的刑法规定的性同意年龄是10是那现在部分的上调到16岁,有人可能也会说老师我印象中的刑事责任年龄,好像这次也是有条件的下调到12岁。',
            '吧唉还真是我记得当时就有同学问我,他说老师你看你主张性的同一年龄上调,同时你又主张刑事责任年龄下调,你不是自我矛盾吗?', '双标。',
            '这些看似矛盾的东西其实不矛盾,因为它针对不同的事物,嘛一个是被动的性防卫联盟,一个是主动的犯罪攻击能力,两者明显不一样,嘛你是张三他爹,同时你又是张一的儿子,说你又是父亲又是儿子,这矛盾与双标吗?',
            '人类思维的复杂性就在于对不同的事物要适用不同的标准,张三是张四的父亲,大黄狗是小黄狗的父亲,张三是父亲,大黄狗也是父亲,所以张三等于大黄狗。',
            '这是苏格拉底当初批评诡辩学派的逻辑,但是现在好像很多人都喜欢这种诡辩学派的风格,呀这个大家是要警惕的。',
            '所以大家就会发现富有照护职责人员性薪水,它显然是一种缓和的家长主义立法,对未成年人法律要像家长一样,通过限制他们的自由来保护他们。',
            '少女也许身体发育成熟,但他的性的心理发育还没有成熟,所以这部分人就很容易成为成年男性剥削的对象,尤其当成年男性属于老师兼顾人这些对少女有一定优势地位或处于信任地位的人,就很容易滥用这种信任地位来剥削少女的心理。',
            '所以刑法规定的这个罪名总之如果是老师,无论是奸淫还是猥亵女生,即便女生同意,16岁以下也可以构成这个罪。',
            '那第三是茂名顶替者,相信大家还记得去年许多茂名顶替的热点案例,那法律就回应了民众的呼声,增加了这个犯罪。',
            '如果盗用冒用他人的身份,顶替他人取得这个高等学历教学入学资格,公务员录用资格,就业安置待遇最高可以出三注意。', '只有上大学考公务员就业安置,冒名顶替才构成这个犯罪,其余地方的冒名顶替那就不构成这个罪。',
            '当然了有可能构成欺诈犯罪。', '前天我同学问我,所以我特别讨厌一个老师,所以我准备每次做完坏事就说是我老师做的。', '这犯罪吗?',
            '比如说张三去特殊娱乐场所,最后对性工作者说我是某某学校教刑法的罗老师,那这肯定不构成冒名顶替者,但也许构成诽谤罪。',
            '当然了,如果张三看到孩子落水,然后把孩子救上岸,很多人问他尊姓大名,张三说我做好事不留名。',
            '那后来实在没办法,就是说我其实姓罗,是逼上一个讲刑法的二不,那这种冒名顶替应该是没事的,当然修正案还增加了其他的一些犯罪,比如说破坏自然保护地址,你在秦岭大兴土木,违反自然保护地的管理法规,建大别墅就可以构成这个罪。',
            '犹如非法植入基因编辑克隆胚胎罪,相信大家还记得当初的贺建奎编辑基因的行为,当初就不太好定罪,后来定了非法行医罪,但其实是比较牵强的。',
            '那现在有了这些罪名,就可以名正言顺的定罪量刑,大家就会发现刑法的很多修改都是对于既定的社会事事务的一个回馈,甘地说能够毁灭人类的有7种食物,其中一种就是没有人性的科学,它有一种就是没有是非的知识,其实修正案还规定了一些新的轻者,但其实都属于现说刑罚权的。',
            '因为这些轻罪以前都是按照什么寻衅滋事罪,以危险方法危害公共安全罪等模糊性罪名处理的,他刑罚是偏重的,而且模糊性罪名的尺度也不太统一,就容易导致这种选择性执法。',
            '你比如说高空抛物坠妨碍安全驾驶罪,这以前经常是按照以危险方法危害公共安全罪处理,但是以危险方法危害公共安全罪的刑罚是很重的。',
            '它的基本形式3~10年,但现在高空抛物罪和妨碍安全驾驶罪的最高刑是一年,这其实是现说了刑罚权避免打击过度,那因此如果高空抛物或者妨碍安全驾驶,危机公共安全最多只能判一年。',
            '但如果严重的危及公共安全,造成了人员伤亡,那还是可以构成以危险方法危害公共安全罪。',
            '就像我们前次刚刚讲过,上从56扔了80个燃烧的蜂窝煤,砸死砸伤了50个人,那就妥妥的以危险方法危害公共安全罪的结果加重犯,最高是可以判处死刑的,那又如催收非法债务罪,侵害英雄烈士名誉荣誉罪,最高刑都是三年,这些行为以前都是可能以寻衅滋事罪论处的。',
            '寻衅滋事罪的基本型是5年以下,加重情节是可以判到10年的。', '所以这两个罪名有不少学者认为是收缩了刑法权,而不是扩张了刑法权。',
            '当然了在修正案中还有一种现象,就是既不是将无罪变为有罪的新增罪名,也不是收缩刑罚权的创设轻罪,而是将以前犯罪的部分类型变更为更重的罪名。',
            '你比如说洗颈椎洗颈椎以前是妨碍公务罪的从重情节,吧最高是三年,那现在规定为一个独立的罪名。', '当它基本性还是三年一下暴力袭击正在依法执行职务的人民警察处,三年以下有期徒刑拘役或者管制。',
            '但是有一个新的加重情形,如果使用枪支管制刀具或者驾驶机动车撞击等手段,严重危及警察的人身安全是可以处3~7年的,只不过这个罪名的对象是警察,如果砸的是警车,让警察感到恐惧,那其实不构成其定罪,那还是可以考虑定以前的妨害公务罪。',
            '今天笼统的给大家讲了讲修正案的新罪名以后,我们再详细介绍,当然了遏制犯罪刑法只是一种最后手段,不到万不得已不应该轻易使用。',
            '我们每个人心中都有黑暗的成分,最重要是约束我们内心的黑暗,让心中的黑暗不至于外泄,成为黑液液体外流,成为自己和社会的责任。', '谢谢各位。']

