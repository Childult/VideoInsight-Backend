package server

import (
	"context"
	"path/filepath"
	"swc/data/absvideo"
	"swc/data/job"
	"swc/data/resource"
	"swc/logger"
	pb "swc/server/network"
	"swc/util"

	"google.golang.org/grpc"
)

func extractVideoAbstract(job *job.Job) {
	logger.Info.Printf("[提取视频摘要] 开始: %+v.\n", job)
	// 判断资源是否已经存在
	av := absvideo.AbsVideo{URL: job.URL}

	if av.ExistInMongodb() { // 视频摘要已经完成
		if job.Status&util.JobTextAbstractExtractionDone != 0 { // 如果文本摘要也已经完成
			job.Status = util.JobCompleted
			job.Save()
			go JobSchedule(job)
		} else {
			job.Status = job.Status | util.JobVideoAbstractExtractionDone
			job.Save()
		}
	} else {
		// 不存在就进行视频分析, 否则忽略
		go videoAnalysis(job)
	}
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
	r := resource.Resource{URL: job.URL}
	err := r.Retrieve()
	if err != nil {
		// 获取资源出错
		logger.Error.Printf("[视频分析] 获取资源出错: %+v.\n", err)
		job.Status = util.JobErrFailedToFindResource
		job.Save()
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
		job.Status = util.JobErrVideoAnalysisGRPCConnectFailed
		job.Save()
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
		job.Status = util.JobErrVideoAnalysisGRPCallFailed
		job.Save()
		go JobSchedule(job)
		return
	}
	logger.Debug.Printf("[视频分析] gRPC 结果: %+v.\n", rpcResult)

	// 结果校验
	jobid := rpcResult.GetJobID()
	if jobid != job.JobID {
		logger.Error.Println("[视频分析] JobID 不匹配")
		job.Status = util.JobErrVideoAnalysisGRPCallFailed
		job.Save()
		go JobSchedule(job)
		return
	}

	// 提取结果
	var videoPath videoAbstract
	videoPath.VAbstract = rpcResult.GetPicName()
	videoPath.Error = rpcResult.GetError()
	if videoPath.Error != "" {
		logger.Error.Printf("[视频分析] 视频分析失败: %+v\n", videoPath.Error)
		job.Status = util.JobErrVideoAnalysisFailed
		job.Save()
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
	absvideo.Dump()
	if job.Status&util.JobTextAbstractExtractionDone != 0 { // 如果文本摘要已经完成
		job.Status = util.JobCompleted // 该任务完成
		job.Save()
	} else {
		job.Status = util.JobVideoAbstractExtractionDone // 文本摘要未完成, 则只将视频摘要的位置置为1, 返回
		job.Save()
	}
	go JobSchedule(job)
}
