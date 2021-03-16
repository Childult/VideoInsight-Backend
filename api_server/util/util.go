package util

import (
	"crypto/sha1"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// 资源定位
const (
	LogFile     = "/swc/log/"          // 日志保存位置
	WorkSpace   = "/swc/code/"         // 工作目录, 用于访问其他文件(如python)
	Location    = "/swc/resource"      // 资源存储位置
	GRPCAddress = "192.168.2.80:50051" // gRPC 调用地址
)

// 任务状态值
const (
	JobStart                       = 1 << iota // 创建资源, 写入数据库
	JobDownloadMedia                           // 准备下载资源
	JobExisted                                 // 文件已存在
	JobExtractAudio                            // 准备提取音频
	JobExtractAudioDone                        // 音频提取成功
	JobTextAbstractExtractionDone              // 文本摘要提取完成
	JobVideoAbstractExtractionDone             // 视频摘要提取完成
	JobCompleted                               // 完成

	JobErrFailedToFindResource              // 从数据库中读取时发生错误
	JobErrDownloadFailed                    // 资源下载失败
	JobErrExtractFailed                     // 音频提取失败
	JobErrTextAnalysisFailed                // 文本分析失败
	JobErrTextAnalysisReadJSONFailed        // 从文本分析结果中获取JSON失败
	JobErrVideoAnalysisFailed               // 视频分析失败
	JobErrVideoAnalysisReadJSONFailed       // 视频分析JSON读取失败
	JobErrVideoAnalysisGRPCConnectFailed    // 视频分析 gRPC 连接失败
	JobErrVideoAnalysisGRPCallFailed        // 视频分析 gRPC 调用失败
	JobErrVideoAnalysisGRPCallJobIDNotMatch // 视频分析 gRPC 调用 JobID 不匹配
)

func GetJobStatus(status int32) string {
	var s []string

	if status&JobStart != 0 {
		s = append(s, "创建资源, 写入数据库")
	}
	if status&JobDownloadMedia != 0 {
		s = append(s, "准备下载资源")
	}
	if status&JobExisted != 0 {
		s = append(s, "文件已存在")
	}
	if status&JobExtractAudio != 0 {
		s = append(s, "准备提取音频")
	}
	if status&JobExtractAudioDone != 0 {
		s = append(s, "音频提取成功")
	}
	if status&JobTextAbstractExtractionDone != 0 {
		s = append(s, "文本摘要提取完成")
	}
	if status&JobVideoAbstractExtractionDone != 0 {
		s = append(s, "视频摘要提取完成")
	}
	if status&JobCompleted != 0 {
		s = append(s, "完成")
	}
	if status&JobErrFailedToFindResource != 0 {
		s = append(s, "从数据库中读取时发生错误")
	}
	if status&JobErrDownloadFailed != 0 {
		s = append(s, "资源下载失败")
	}
	if status&JobErrExtractFailed != 0 {
		s = append(s, "音频提取失败")
	}
	if status&JobErrTextAnalysisFailed != 0 {
		s = append(s, "文本分析失败")
	}
	if status&JobErrTextAnalysisReadJSONFailed != 0 {
		s = append(s, "从文本分析结果中获取JSON失败")
	}
	if status&JobErrVideoAnalysisFailed != 0 {
		s = append(s, "视频分析失败")
	}
	if status&JobErrVideoAnalysisReadJSONFailed != 0 {
		s = append(s, "视频分析JSON读取失败")
	}
	if status&JobErrVideoAnalysisGRPCConnectFailed != 0 {
		s = append(s, "视频分析 gRPC 连接失败")
	}
	if status&JobErrVideoAnalysisGRPCallFailed != 0 {
		s = append(s, "视频分析 gRPC 调用失败")
	}
	if status&JobErrVideoAnalysisGRPCallJobIDNotMatch != 0 {
		s = append(s, "视频分析 gRPC 调用 JobID 不匹配")
	}
	return strings.Join(s, "->")
}

// 资源状态值
const (
	ResourceDownloading = 1 << iota // 下载资源
	ResourceExtracting              // 提取音频
	ResourceCompleted               // 完成

	ResourceErrDownloadFailed // 资源下载失败
	ResourceErrExtractFailed  // 音频提取失败
)

// MessageJSON 用户利用 POST 提交的数据, 用于为任务创建唯一的 ID
type MessageJSON struct {
	DeviceID string   `json:"device_id"`
	URL      string   `json:"url"`
	KeyWords []string `json:"keywords,omitempty"`
}

// String toString, 用于构建 hash, 最终返回唯一ID
func (json MessageJSON) String() string {
	return json.DeviceID + json.URL + strings.Join(json.KeyWords, "")
}

// GetHash 返回固定大小的 hash 值
func (json MessageJSON) GetHash() (result [12]byte) {
	hash := sha1.New()
	hash.Write([]byte(json.String()))
	copy(result[:], hash.Sum([]byte(""))[0:12])
	return
}

// GetID 通过哈希返回唯一 ID
func (json MessageJSON) GetID() string {
	return fmt.Sprintf("%x", json.GetHash())
}

// GetJSON 从用户输入中构建 MessageJSON, KeyWords 为空时设为空切片 []string{}
func GetJSON(c *gin.Context) (json MessageJSON, err error) {
	// 获取数据
	err = c.ShouldBindJSON(&json)
	if err != nil {
		err = fmt.Errorf("格式错误.")
		return
	}
	json.KeyWords = removeEmptyString(json.KeyWords)
	return
}

// removeEmptyString 删除切片中的空串
func removeEmptyString(a []string) []string {
	return deleteKeywords(a, "")
}

// deleteKeywords 删除切片中指定字符串, 并且希望原始切片为 nil 时, 返回一个空的切片 []string{}
func deleteKeywords(rawSlice []string, target string) []string {
	len := len(rawSlice)
	newSlice := make([]string, len)
	i := 0
	for j := 0; j < len; j++ {
		if rawSlice[j] != target {
			newSlice[i] = rawSlice[j]
			i++
		}
	}
	return newSlice[:i]
}
