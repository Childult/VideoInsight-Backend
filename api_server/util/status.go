package util

// 资源定位
const (
	LogFile   = "/swc/log/"          // 日志保存位置
	WorkSpace = "/swc/code/"         // 工作目录, 用于访问其他文件(如python)
	Location  = "/swc/resource"      // 资源存储位置
	MongoAddr = "192.168.2.80:27018" // Mongodb 地址
	MongoUser = ""                   // Mongodb 账户
	MongoPW   = ""                   // Mongodb 密码
	RedisAddr = "192.168.2.80:6379"  // Redis 地址
	RedisPW   = ""                   // Redis 密码
)

var (
	MongoDB     = "swc"                // mongodb 数据库名称
	RedisDB     = 0                    // redis 数据库号数
	GRPCAddress = "192.168.2.80:50051" // gRPC 调用地址
)

// 任务状态值
const (
	TaskStart        = iota // 接收到一个任务 创建资源, 写入数据库
	TaskCreated             // 任务创建
	TaskDownloadDone        // 资源下载完成
	TaskCompleted           // 完成

	TaskErrRetrieveFail        // 资源获取出错
	TaskErrTextAnalysisFailed  // 文本分析失败
	TaskErrVideoAnalysisFailed // 视频分析失败
	TaskErrTextAndVideoFail    // 文本分析和视频分析都失败
)

// 资源状态值
const (
	ResourceCreated         = iota // 资源创建
	ResourceDownloading            // 下载资源
	ResourceDownloadDone           // 下载成功
	ResourceAudioExtracting        // 提取音频
	ResourceCompleted              // 完成

	ResourceErrDownloadFailed     // 资源下载失败
	ResourceErrAudioExtractFailed // 音频提取失败
)

const (
	AbsTextStart                         = iota // 创建资源
	AbsTextComplete                             // 文本分析完成
	AbsTextErrPythonFailed                      // 无返回结果, python 调用出错
	AbsTextErrTextAnalysisReadJSONFailed        // 从文本分析结果中获取JSON失败
	AbsTextErrTextAnalysisFailed                // 文本分析失败
)

const (
	AbsVideoStart                   = iota // 创建视频摘要资源
	AbsVideoComplete                       // 视频分析完成
	AbsVideoErrGRPCConnectFailed           // 视频分析 gRPC 连接失败
	AbsVideoErrGRPCallFailed               // 视频分析 gRPC 调用失败
	AbsVideoErrGRPCallJobIDNotMatch        // 视频分析 gRPC 调用 JobID 不匹配
)
