# -*- encoding:utf-8 -*-
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
    d = {'data': '[{"bg":"20","ed":"1010","onebest":"嗯","speaker":"0"},{"bg":"5550","ed":"12970","onebest":"各位同学大家好，啊对于刑法学的唉呀今年3月1号呢是一个大日子，因为刑法修正案11啊生效了，","speaker":"0"},{"bg":"12990","ed":"15550","onebest":"刑法中呢增加了许多新的罪名，","speaker":"0"},{"bg":"16050","ed":"17980","onebest":"法网呢被进一步的严密。","speaker":"0"},{"bg":"18710","ed":"23960","onebest":"所以今天啊想给大家简单介绍一下，修正案11呢主要增加了哪些新罪名，","speaker":"0"},{"bg":"24400","ed":"27360","onebest":"哪些行为已经开始变成犯罪了？","speaker":"0"},{"bg":"27550","ed":"35390","onebest":"呃我先想问大家一个问题，大家觉得刑法修正案11是扩张的刑罚权还是收缩了刑罚权？","speaker":"0"},{"bg":"35390","ed":"35750","onebest":"呢","speaker":"0"},{"bg":"35790","ed":"37220","onebest":"扩张咱选一，","speaker":"0"},{"bg":"37880","ed":"40290","onebest":"收缩咱们选啊答案，呢","speaker":"0"},{"bg":"41130","ed":"42310","onebest":"当然是老规矩，","speaker":"0"},{"bg":"43020","ed":"48880","onebest":"应该选商刑法修正案，啊其实系扩张的刑罚权有收缩的刑罚权，","speaker":"0"},{"bg":"49000","ed":"52570","onebest":"所以这是有些新增的罪名啊属于真正的新增。","speaker":"0"},{"bg":"52750","ed":"57770","onebest":"也就是说以前这些行为不构成犯罪，3月1号之后才构成犯罪的。","speaker":"0"},{"bg":"57980","ed":"59840","onebest":"那根据从旧兼从轻原则，","speaker":"0"},{"bg":"60030","ed":"66220","onebest":"如果一种新的法律对行为人不利，那就不能溯及既往，毕竟刑法嘛是一根带哨子的鞭子，","speaker":"0"},{"bg":"66640","ed":"68650","onebest":"在抽你之前是要告诉你的，","speaker":"0"},{"bg":"68920","ed":"74330","onebest":"那么这些让无罪变为有罪的真正意义上的新尊罪名有哪些呢？","speaker":"0"},{"bg":"74990","ed":"78260","onebest":"那主要有如下一些，啊一个呢是商业间谍行为，","speaker":"0"},{"bg":"78670","ed":"84250","onebest":"为境外的机构组织人员窃取、刺探、收买非法提供商业秘密的，","speaker":"0"},{"bg":"84830","ed":"92180","onebest":"啊但是可以处5年以下有期徒刑，那么情节特别严重的，甚至还可以触到5年以上，最高可以判到15年有期徒刑。","speaker":"0"},{"bg":"92280","ed":"93620","onebest":"因为刑法以前啊","speaker":"0"},{"bg":"93900","ed":"97970","onebest":"只有为境外窃取刺探收买非法提供国家秘密，","speaker":"0"},{"bg":"98520","ed":"103740","onebest":"但现在商业间谍呢是越来越多，所以就增加了这个犯罪比茅台酒的配方啊","speaker":"0"},{"bg":"104090","ed":"110110","onebest":"加工的独特方法。虽然茅台酒厂是国企，但这个肯定不属于国家秘密，嘛","speaker":"0"},{"bg":"110530","ed":"111680","onebest":"但它属于商业秘密。","speaker":"0"},{"bg":"112160","ed":"116320","onebest":"那如果张三把这个独特的加工方法告诉了外国的酒厂，","speaker":"0"},{"bg":"116760","ed":"118030","onebest":"那就有可能构成。","speaker":"0"},{"bg":"118330","ed":"119240","onebest":"那第二呢","speaker":"0"},{"bg":"119430","ed":"123840","onebest":"是负有照护职责，人员性侵罪这个我们以前也提过，","speaker":"0"},{"bg":"124190","ed":"131370","onebest":"对已满14周岁不满16周岁的未成年女性负有监护、啊收养、啊看护、教育、医疗","speaker":"0"},{"bg":"131640","ed":"133110","onebest":"等特殊职责的人员，","speaker":"0"},{"bg":"133450","ed":"137080","onebest":"与该未成年女性发生性关系处三年以下有期徒刑，","speaker":"0"},{"bg":"137430","ed":"138650","onebest":"那情节恶劣，呢","speaker":"0"},{"bg":"138830","ed":"140140","onebest":"可以出3~10年。","speaker":"0"},{"bg":"140710","ed":"144490","onebest":"那大家知道以前的刑法规定的性同意年龄是10是","speaker":"0"},{"bg":"144910","ed":"147170","onebest":"那现在部分的上调到16岁，","speaker":"0"},{"bg":"147780","ed":"154460","onebest":"有人可能也会说老师我印象中的刑事责任年龄，好像这次也是有条件的下调到12岁。吧","speaker":"0"},{"bg":"154720","ed":"157190","onebest":"唉还真是我记得当时就有同学问我，","speaker":"0"},{"bg":"157410","ed":"161380","onebest":"他说老师你看你主张性的同一年龄上调，","speaker":"0"},{"bg":"161690","ed":"165360","onebest":"同时你又主张刑事责任年龄下调，","speaker":"0"},{"bg":"165980","ed":"168180","onebest":"你不是自我矛盾吗？双标。","speaker":"0"},{"bg":"168180","ed":"168550","onebest":"啊","speaker":"0"},{"bg":"168750","ed":"171510","onebest":"这些看似矛盾的东西啊其实不矛盾，","speaker":"0"},{"bg":"172130","ed":"173940","onebest":"因为它针对不同的事物，嘛","speaker":"0"},{"bg":"174350","ed":"176590","onebest":"一个是被动的性防卫联盟，","speaker":"0"},{"bg":"176980","ed":"179430","onebest":"一个是主动的犯罪攻击能力，","speaker":"0"},{"bg":"179880","ed":"185970","onebest":"两者明显不一样，嘛你是张三他爹，同时你又是张一的儿子，","speaker":"0"},{"bg":"186440","ed":"189260","onebest":"说你又是父亲又是儿子，","speaker":"0"},{"bg":"189570","ed":"191010","onebest":"这矛盾与双标吗？","speaker":"0"},{"bg":"191010","ed":"192960","onebest":"人类思维的复杂性","speaker":"0"},{"bg":"193330","ed":"198650","onebest":"就在于对不同的事物要适用不同的标准，张三是张四的父亲，","speaker":"0"},{"bg":"198830","ed":"201070","onebest":"大黄狗是小黄狗的父亲，","speaker":"0"},{"bg":"201600","ed":"202690","onebest":"张三是父亲，","speaker":"0"},{"bg":"203130","ed":"206710","onebest":"大黄狗也是父亲，所以张三等于大黄狗。","speaker":"0"},{"bg":"206910","ed":"214970","onebest":"这是苏格拉底当初批评诡辩学派的逻辑，但是现在好像很多人都喜欢这种诡辩学派的风格，呀","speaker":"0"},{"bg":"215500","ed":"216770","onebest":"这个大家是要警惕的。","speaker":"0"},{"bg":"217020","ed":"220960","onebest":"所以大家就会发现富有照护职责人员性薪水，","speaker":"0"},{"bg":"221340","ed":"227940","onebest":"它显然是一种缓和的家长主义立法，对未成年人法律要像家长一样，通过限制他们的自由来保护他们。","speaker":"0"},{"bg":"228010","ed":"230060","onebest":"少女也许身体发育成熟，","speaker":"0"},{"bg":"230440","ed":"233050","onebest":"但他的性的心理发育还没有成熟，","speaker":"0"},{"bg":"233470","ed":"241650","onebest":"所以这部分人就很容易成为成年男性剥削的对象，尤其当成年男性属于老师兼顾人这些对少女","speaker":"0"},{"bg":"241890","ed":"244750","onebest":"啊有一定优势地位或处于信任地位的人，","speaker":"0"},{"bg":"244980","ed":"248940","onebest":"就很容易滥用这种信任地位来剥削少女的心理。","speaker":"0"},{"bg":"249450","ed":"251020","onebest":"所以刑法规定的这个罪名","speaker":"0"},{"bg":"251700","ed":"253080","onebest":"啊总之如果是老师，","speaker":"0"},{"bg":"253310","ed":"256620","onebest":"无论是奸淫还是猥亵女生，即便女生同意，","speaker":"0"},{"bg":"256820","ed":"259020","onebest":"啊16岁以下也可以构成这个罪。","speaker":"0"},{"bg":"259190","ed":"266350","onebest":"那第三呢是茂名顶替者，相信大家还记得去年许多茂名顶替的热点案例，啊那法律就回应了民众的呼声，","speaker":"0"},{"bg":"266390","ed":"267380","onebest":"增加了这个犯罪。","speaker":"0"},{"bg":"267560","ed":"273540","onebest":"如果盗用冒用他人的身份，顶替他人取得这个高等学历教学入学资格，","speaker":"0"},{"bg":"273870","ed":"277840","onebest":"公务员录用资格，就业安置待遇最高可以出三注意。啊","speaker":"0"},{"bg":"278110","ed":"283130","onebest":"只有上大学考公务员就业安置，冒名顶替才构成这个犯罪，","speaker":"0"},{"bg":"283570","ed":"286220","onebest":"其余地方的冒名顶替那就不构成这个罪。","speaker":"0"},{"bg":"286540","ed":"288210","onebest":"当然了有可能构成欺诈犯罪。啊","speaker":"0"},{"bg":"288390","ed":"294830","onebest":"前天我同学问我，所以我特别讨厌一个老师，所以我准备每次做完坏事就说是我老师做的。","speaker":"0"},{"bg":"294870","ed":"297490","onebest":"这犯罪吗？比如说张三去特殊娱乐场所，","speaker":"0"},{"bg":"298140","ed":"299810","onebest":"最后对性工作者说我","speaker":"0"},{"bg":"300120","ed":"304720","onebest":"是某某学校教刑法的罗老师，那这肯定不构成冒名顶替者，","speaker":"0"},{"bg":"305150","ed":"306370","onebest":"但也许构成诽谤罪。","speaker":"0"},{"bg":"306550","ed":"306980","onebest":"当然了，","speaker":"0"},{"bg":"307180","ed":"310480","onebest":"如果张三看到孩子落水，然后把孩子救上岸，","speaker":"0"},{"bg":"310960","ed":"314120","onebest":"很多人问他尊姓大名，张三说我做好事不留名。","speaker":"0"},{"bg":"314390","ed":"319560","onebest":"那后来实在没办法，啊就是说我其实姓罗，啊是逼上一个讲刑法的二不，","speaker":"0"},{"bg":"319630","ed":"326590","onebest":"那这种冒名顶替呢应该是没事的，当然修正案还增加了其他的一些犯罪，比如说破坏自然保护地址，","speaker":"0"},{"bg":"327000","ed":"328640","onebest":"啊你在秦岭大兴土木，","speaker":"0"},{"bg":"328970","ed":"333210","onebest":"违反自然保护地的管理法规，建大别墅就可以构成这个罪。","speaker":"0"},{"bg":"333420","ed":"340040","onebest":"犹如非法植入基因编辑克隆胚胎罪，相信大家还记得当初的贺建奎啊编辑基因的行为，","speaker":"0"},{"bg":"340520","ed":"345120","onebest":"当初就不太好定罪，后来定了非法行医罪，但其实是比较牵强的。","speaker":"0"},{"bg":"345330","ed":"346920","onebest":"那现在有了这些罪名，","speaker":"0"},{"bg":"347150","ed":"348910","onebest":"就可以名正言顺的","speaker":"0"},{"bg":"349110","ed":"352670","onebest":"定罪量刑，大家就会发现刑法的很多修改","speaker":"0"},{"bg":"353010","ed":"355430","onebest":"都是对于既定的社会事事务","speaker":"0"},{"bg":"355630","ed":"359230","onebest":"的一个回馈，甘地说能够毁灭人类的有7种食物，","speaker":"0"},{"bg":"359630","ed":"363030","onebest":"其中一种就是没有人性的科学，它有一种","speaker":"0"},{"bg":"363590","ed":"368940","onebest":"就是没有是非的知识，其实呢修正案还规定了一些新的轻者，","speaker":"0"},{"bg":"369160","ed":"371540","onebest":"但其实都属于现说刑罚权的。","speaker":"0"},{"bg":"371720","ed":"378010","onebest":"因为这些轻罪啊以前都是按照什么寻衅滋事罪，啊以危险方法危害公共安全罪","speaker":"0"},{"bg":"378030","ed":"379590","onebest":"等模糊性罪名处理的，","speaker":"0"},{"bg":"380090","ed":"381330","onebest":"他刑罚是偏重的，","speaker":"0"},{"bg":"381970","ed":"385120","onebest":"而且模糊性罪名的尺度啊也不太统一，","speaker":"0"},{"bg":"385410","ed":"387850","onebest":"就容易导致这种选择性执法。","speaker":"0"},{"bg":"388030","ed":"390920","onebest":"你比如说高空抛物坠妨碍安全驾驶罪，","speaker":"0"},{"bg":"391390","ed":"394060","onebest":"这以前经常是按照以危险方法","speaker":"0"},{"bg":"394250","ed":"396230","onebest":"危害公共安全罪处理，","speaker":"0"},{"bg":"396760","ed":"399820","onebest":"但是以危险方法危害公共安全罪的刑罚是很重的。","speaker":"0"},{"bg":"400220","ed":"401910","onebest":"它的基本形式3~10年，","speaker":"0"},{"bg":"402390","ed":"406300","onebest":"但现在高空抛物罪和妨碍安全驾驶罪的最高刑是一年，","speaker":"0"},{"bg":"406660","ed":"412770","onebest":"这其实是现说了刑罚权避免打击过度，那因此如果高空抛物或者妨碍安全驾驶，","speaker":"0"},{"bg":"412790","ed":"414860","onebest":"啊危机公共安全最多只能判一年。","speaker":"0"},{"bg":"415260","ed":"418900","onebest":"但如果严重的危及公共安全，造成了人员伤亡，","speaker":"0"},{"bg":"419250","ed":"421900","onebest":"那还是可以构成以危险方法危害公共安全罪。","speaker":"0"},{"bg":"422600","ed":"424540","onebest":"就像我们前次刚刚讲过，啊","speaker":"0"},{"bg":"424710","ed":"427910","onebest":"上从56扔了80个燃烧的蜂窝煤，","speaker":"0"},{"bg":"428140","ed":"430070","onebest":"啊砸死砸伤了50个人，","speaker":"0"},{"bg":"430370","ed":"433860","onebest":"那就妥妥的以危险方法危害公共安全罪的结果加重犯，","speaker":"0"},{"bg":"434210","ed":"437890","onebest":"最高是可以判处死刑的，那又如催收非法债务罪，啊","speaker":"0"},{"bg":"438190","ed":"447030","onebest":"侵害英雄烈士名誉荣誉罪，最高刑都是三年，这些行为啊以前呢都是可能以寻衅滋事罪论处的。","speaker":"0"},{"bg":"447110","ed":"450020","onebest":"啊寻衅滋事罪的基本型是5年以下，","speaker":"0"},{"bg":"450380","ed":"452730","onebest":"加重情节是可以判到10年的。","speaker":"0"},{"bg":"453100","ed":"459010","onebest":"所以这两个罪名啊有不少学者认为是收缩了刑法权，而不是扩张了刑法权。","speaker":"0"},{"bg":"459040","ed":"465790","onebest":"当然了在修正案中啊还有一种现象，就是既不是将无罪变为有罪的新增罪名，","speaker":"0"},{"bg":"465970","ed":"469160","onebest":"也不是收缩刑罚权的创设轻罪，","speaker":"0"},{"bg":"469670","ed":"473670","onebest":"而是将以前犯罪的部分类型变更为更重的罪名。","speaker":"0"},{"bg":"474010","ed":"477650","onebest":"你比如说洗颈椎洗颈椎以前是妨碍公务罪的从重情节，吧","speaker":"0"},{"bg":"477830","ed":"478610","onebest":"最高是三年，","speaker":"0"},{"bg":"478970","ed":"481240","onebest":"那现在规定为一个独立的罪名。","speaker":"0"},{"bg":"481460","ed":"488320","onebest":"当它基本性还是三年一下暴力袭击正在依法执行职务的人民警察处，三年以下有期徒刑拘役或者管制。","speaker":"0"},{"bg":"488770","ed":"495170","onebest":"但是有一个新的加重情形，如果使用枪支啊管制刀具啊或者驾驶机动车撞击等手段，","speaker":"0"},{"bg":"495280","ed":"496370","onebest":"严重危及","speaker":"0"},{"bg":"496590","ed":"499450","onebest":"啊警察的人身安全是可以处3~7年的，","speaker":"0"},{"bg":"499630","ed":"503170","onebest":"只不过这个罪名的对象是警察，如果砸的是警车，","speaker":"0"},{"bg":"503430","ed":"509720","onebest":"让警察感到恐惧，那其实不构成其定罪，那还是可以考虑定以前的妨害公务罪。","speaker":"0"},{"bg":"509750","ed":"515400","onebest":"呃今天呢笼统的给大家讲了讲修正案的新罪名啊以后，呢我们再详细介绍，当然了","speaker":"0"},{"bg":"515680","ed":"518030","onebest":"遏制犯罪啊刑法只是一种最后手段，","speaker":"0"},{"bg":"518530","ed":"520740","onebest":"不到万不得已不应该轻易使用。","speaker":"0"},{"bg":"520910","ed":"523030","onebest":"我们每个人心中都有黑暗的成分，","speaker":"0"},{"bg":"523320","ed":"523930","onebest":"最重要","speaker":"0"},{"bg":"524110","ed":"528400","onebest":"是约束我们内心的黑暗，让心中的黑暗不至于外泄，","speaker":"0"},{"bg":"529150","ed":"531310","onebest":"成为黑液液体外流，","speaker":"0"},{"bg":"531720","ed":"532520","onebest":"成为自己","speaker":"0"},{"bg":"532710","ed":"533400","onebest":"和社会的责任。","speaker":"0"},{"bg":"533910","ed":"534740","onebest":"谢谢各位。","speaker":"0"}]', 'err_no': 0, 'failed': None, 'ok': 0}
    print(preprocess_audio_text(d))

    # t = extract_document('AText.json')
    # for seg in sen_generator(t):
    #     print(len(seg), seg)
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
