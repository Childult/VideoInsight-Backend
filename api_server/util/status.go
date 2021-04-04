package util

// 资源定位
const (
	LogFile     = "/swc/log/"          // 日志保存位置
	WorkSpace   = "/swc/code/"         // 工作目录, 用于访问其他文件(如python)
	Location    = "/swc/resource"      // 资源存储位置
	GRPCAddress = "192.168.2.80:50051" // gRPC 调用地址
	MongoAddr   = "192.168.2.80:27018" // Mongodb 地址
	MongoUser   = ""                   // Mongodb 账户
	MongoPW     = ""                   // Mongodb 密码
	RedisAddr   = "192.168.2.80:6379"  // Redis 地址
	RedisPW     = ""                   // Redis 密码
)

var (
	MongoDB = "swcdb" // mongodb 数据库名称
	RedisDB = 0       // redis 数据库号数
)

// 任务状态值
const (
	JobStart                       = iota // 创建资源, 写入数据库
	JobResourceExisted                    // 资源已存在
	JobToTextExtract                      // 准备提取文本摘要
	JobToDownloadMedia                    // 准备下载资源
	JobToExtractAudio                     // 准备提取音频
	JobAudioExtractDone                   // 音频提取成功
	JobTextAbstractExtractionDone         // 文本摘要提取完成
	JobVideoAbstractExtractionDone        // 视频摘要提取完成
	JobCompleted                          // 完成

	JobErrFailedToFindResource              // 从数据库中读取时发生错误
	JobErrDownloadFailed                    // 资源下载失败
	JobErrAudioExtractFailed                // 音频提取失败
	JobErrTextAnalysisFailed                // 文本分析失败
	JobErrTextAnalysisReadJSONFailed        // 从文本分析结果中获取JSON失败
	JobErrVideoAnalysisFailed               // 视频分析失败
	JobErrVideoAnalysisReadJSONFailed       // 视频分析JSON读取失败
	JobErrVideoAnalysisGRPCConnectFailed    // 视频分析 gRPC 连接失败
	JobErrVideoAnalysisGRPCallFailed        // 视频分析 gRPC 调用失败
	JobErrVideoAnalysisGRPCallJobIDNotMatch // 视频分析 gRPC 调用 JobID 不匹配
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
