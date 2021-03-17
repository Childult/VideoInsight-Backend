package server

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"swc/logger"
	"swc/mongodb"
	"swc/mongodb/abstext"
	"swc/mongodb/absvideo"
	"swc/mongodb/job"
	"swc/mongodb/resource"
	pb "swc/server/network"
	"swc/util"

	"google.golang.org/grpc"
)

// textAnalysis 文本分析
// 不判断该资源是否做过文本分析, 直接插入数据库
func textAnalysis(job *job.Job) {
	logger.Info.Printf("[文本分析] 开始: %+v.\n", job)
	// 获取资源信息
	r, err := resource.GetByKey(job.URL)
	if err != nil {
		logger.Error.Printf("[文本分析] 获取资源出错: %+v.\n", err)
		job.SetStatus(util.JobErrFailedToFindResource)
		go JobSchedule(job)
		return
	}

	// 构建文本分析对象
	python := PyWorker{
		PackagePath: filepath.Join(util.WorkSpace, "text_analysis"), // python 文件所在的包
		FileName:    "api",                                          // 文件名
		MethodName:  "generate_abstract_from_audio",                 // 调用函数名
		Args: []string{ // 调用实参
			SetArg(filepath.Join(r.Location, r.AudioPath)),
		},
	}

	// 文本分析
	go python.Call(job, textHandle)
}

// textAbstract 用于存储文本分析的结果
type textAbstract struct {
	AText     string `json:"AText"`
	TAbstract string `json:"TAbstract"`
	Error     string `json:"Error"`
}

// textHandle 文本分析的回调
func textHandle(job *job.Job, result []string) {
	logger.Info.Printf("[文本分析回调] 开始: %+v.\n", job)
	// 获取资源信息
	r, err := resource.GetByKey(job.URL)
	if err != nil {
		logger.Error.Printf("[文本分析回调] 获取资源出错: %+v.\n", err)
		job.SetStatus(util.JobErrFailedToFindResource)
		go JobSchedule(job)
		return
	}

	if len(result) == 0 { // 未找到结果
		logger.Error.Println("[文本分析回调] 文本分析失败.")
		job.SetStatus(util.JobErrTextAnalysisFailed)
		go JobSchedule(job)
		return
	}
	pythonReturn := strings.Join(result, "")

	// 找到结果, 从返回结果中提取数据
	var text textAbstract
	err = json.Unmarshal([]byte(pythonReturn), &text)
	if err != nil {
		logger.Error.Printf("[文本分析回调] 从文本分析结果中获取JSON失败: %+v.\n", err)
		job.SetStatus(util.JobErrTextAnalysisReadJSONFailed)
		go JobSchedule(job)
		return
	}

	if text.Error != "" {
		logger.Error.Printf("[文本分析回调] 文本分析失败: %+v.\n", text.Error)
		job.SetStatus(util.JobErrTextAnalysisFailed)
		go JobSchedule(job)
		return
	}

	// 文本分析成功, 初始化需要存储的数据
	abstext := abstext.NewAbsText(job.URL, text.AText, text.TAbstract, job.KeyWords)
	r.SetAbsText(abstext.Hash)   // 更新资源对应的文本摘要
	job.SetAbsText(abstext.Hash) // 更新任务对应的文本摘要
	mongodb.InsertOne(abstext)   // 将文本分析结果插入数据库

	if job.Status&util.JobVideoAbstractExtractionDone != 0 { // 视频分析已经完成
		job.SetStatus(util.JobCompleted) // 任务结束
	} else {
		job.SetStatus(util.JobTextAbstractExtractionDone) // 将文本分析完成的位置为1
	}
	go JobSchedule(job)
}

// videoAbstract 用于存储视频摘要的结果
type videoAbstract struct {
	VAbstract []string `json:"VAbstract"`
	Error     string   `json:"Error"`
}

// videoAnalysis 视频分析
// 不判断该资源是否做过视频分析, 直接插入数据库
func videoAnalysis(job *job.Job) {
	logger.Info.Printf("[视频分析] 开始: %+v.\n", job)
	// 获取视频资源信息
	r, err := resource.GetByKey(job.URL)
	if err != nil {
		// 获取资源出错
		logger.Error.Printf("[视频分析] 获取资源出错: %+v.\n", err)
		job.SetStatus(util.JobErrFailedToFindResource)
		go JobSchedule(job)
		return
	}

	address := util.GRPCAddress                         // gRPC 地址
	jobID := job.JobID                                  // 任务id
	videoFile := filepath.Join(r.Location, r.VideoPath) // 视频文件路径
	savaPath := r.Location                              // 结果存储路径
	logger.Debug.Printf("[视频分析] gRPC 参数: address:%+v, jobID:%+v, videoFile:%+v, savaPath:%+v.\n", address, jobID, videoFile, savaPath)

	// 连接 gRPC 服务器
	logger.Debug.Println("[视频分析] 连接 gRPC 服务器.")
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.Error.Printf("[视频分析] gRPC 服务器连接失败: %+v.\n", err)
		job.SetStatus(util.JobErrVideoAnalysisGRPCConnectFailed)
		go JobSchedule(job)
		return
	}

	defer conn.Close()                     // 关闭连接
	rpc := pb.NewVideoAnalysisClient(conn) // gRPC 调用的句柄

	// 开始调用服务
	logger.Debug.Println("[视频分析] 开始调用 gRPC 服务.")
	rpcResult, err := rpc.GetStaticVideoAbstract(context.TODO(), &pb.VideoInfo{JobId: jobID, File: videoFile, SaveDir: savaPath})
	if err != nil {
		logger.Error.Printf("[视频分析] gRPC 调用失败: %+v.\n", err)
		job.SetStatus(util.JobErrVideoAnalysisGRPCallFailed)
		go JobSchedule(job)
		return
	}
	logger.Debug.Printf("[视频分析] gRPC 结果: %+v.\n", rpcResult)

	// 结果校验
	jobid := rpcResult.GetJobID()
	if jobid != job.JobID {
		logger.Error.Println("[视频分析] JobID 不匹配")
		job.SetStatus(util.JobErrVideoAnalysisGRPCallFailed)
		go JobSchedule(job)
		return
	}

	// 提取结果
	var videoPath videoAbstract
	videoPath.VAbstract = rpcResult.GetPicName()
	videoPath.Error = rpcResult.GetError()
	if videoPath.Error != "" {
		logger.Error.Printf("[视频分析] 视频分析失败: %+v\n", videoPath.Error)
		job.SetStatus(util.JobErrVideoAnalysisFailed)
		go JobSchedule(job)
		return
	}
	logger.Debug.Printf("[视频分析] 视频分析结果: %+v.\n", videoPath)

	// 初始化需要存储的数据
	absvideo := absvideo.AbsVideo{
		URL:      job.URL,
		Abstract: videoPath.VAbstract,
	}
	// 把数据插入数据库
	mongodb.InsertOne(absvideo)
	if job.Status&util.JobTextAbstractExtractionDone != 0 { // 如果文本摘要已经完成
		job.SetStatus(util.JobCompleted) // 该任务完成
	} else {
		job.SetStatus(util.JobVideoAbstractExtractionDone) // 文本摘要未完成, 则只将视频摘要的位置置为1, 返回
	}
	go JobSchedule(job)
}
