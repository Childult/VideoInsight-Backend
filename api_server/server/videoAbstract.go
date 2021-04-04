package server

import (
	"context"
	"path/filepath"
	"swc/data/absvideo"
	"swc/data/resource"
	"swc/data/task"
	"swc/dbs/mongodb"
	"swc/dbs/redis"
	"swc/logger"
	pb "swc/server/network"
	"swc/util"

	"google.golang.org/grpc"
)

// videoAbstract 用于存储视频摘要的结果
type videoAbstract struct {
	VAbstract []string `json:"VAbstract"`
	Error     string   `json:"Error"`
}

func extractVideoAbstract(t *task.Task, r *resource.Resource) {
	logger.Info.Printf("[提取视频摘要] 视频分析开始: %+v.\n", t)

	// 设置 gRPC 参数
	address := util.GRPCAddress                         // gRPC 地址
	jobID := t.TaskID                                   // 任务id
	videoFile := filepath.Join(r.Location, r.VideoPath) // 视频文件路径
	savaPath := r.Location                              // 结果存储路径
	logger.Debug.Printf("[提取视频摘要] gRPC 参数: address:%+v, jobID:%+v, videoFile:%+v, savaPath:%+v.\n", address, jobID, videoFile, savaPath)

	// 连接 gRPC 服务器
	logger.Debug.Println("[提取视频摘要] 连接 gRPC 服务器.")
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.Error.Printf("[提取视频摘要] gRPC 服务器连接失败: %+v.\n", err)
		t.Status = util.JobErrVideoAnalysisGRPCConnectFailed
		redis.UpdataOne(t)
		return
	}

	defer conn.Close()                     // 关闭连接
	rpc := pb.NewVideoAnalysisClient(conn) // gRPC 调用的句柄

	// 开始调用服务
	logger.Debug.Println("[提取视频摘要] 开始调用 gRPC 服务.")
	rpcResult, err := rpc.GetStaticVideoAbstract(context.TODO(), &pb.VideoInfo{JobId: jobID, File: videoFile, SaveDir: savaPath})
	if err != nil {
		logger.Error.Printf("[提取视频摘要] gRPC 调用失败: %+v.\n", err)
		t.Status = util.JobErrVideoAnalysisGRPCallFailed
		redis.UpdataOne(t)
		return
	}
	logger.Debug.Printf("[提取视频摘要] gRPC 结果: %+v.\n", rpcResult)

	// 结果校验
	jobid := rpcResult.GetJobID()
	if jobid != t.TaskID {
		logger.Error.Println("[提取视频摘要] JobID 不匹配")
		t.Status = util.JobErrVideoAnalysisGRPCallFailed
		redis.UpdataOne(t)
		return
	}

	// 提取结果
	var videoPath videoAbstract
	videoPath.VAbstract = rpcResult.GetPicName()
	videoPath.Error = rpcResult.GetError()
	if videoPath.Error != "" {
		logger.Error.Printf("[提取视频摘要] 视频分析失败: %+v\n", videoPath.Error)
		t.Status = util.JobErrVideoAnalysisFailed
		redis.UpdataOne(t)
		return
	}
	logger.Debug.Printf("[提取视频摘要] 视频分析结果: %+v.\n", videoPath)

	// 初始化需要存储的数据
	absvideo := absvideo.AbsVideo{
		URL:      t.URL,
		Abstract: videoPath.VAbstract,
	}
	// 把数据插入数据库
	redis.InsertOne(&absvideo)
	mongodb.InsertOne(&absvideo)
	redis.UpdataOne(t)
	setAbstractFlag(t, util.JobVideoAbstractExtractionDone)
}
