package util

import (
	"crypto/sha1"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// 资源定位
const (
	LogFile   = "/swc/log/"     // 日志保存位置
	WorkSpace = "/swc/code/"    // 工作目录, 用于访问其他文件(如python)
	Location  = "/swc/resource" // 资源存储位置
)

// 任务状态值
const (
	JobStart                       = 0  // 创建资源, 写入数据库
	JobDownloadMedia               = 1  // 下载资源
	JobExisted                     = 2  // 文件已存在
	JobExtractAudio                = 3  // 提取音频
	JobExtractAudioDone            = 4  // 音频提取成功
	JobTextAbstractExtractionDone  = 8  // 文本摘要提取完成
	JobVideoAbstractExtractionDone = 16 // 视频摘要提取完成
	JobCompleted                   = 32 // 完成

	JobErrFailedToFindResource        = 100 // 从数据库中读取时发生错误
	JobErrDownloadFailed              = 101 // 资源下载失败
	JobErrExtractFailed               = 102 // 音频提取失败
	JobErrTextAnalysisFailed          = 103 // 文本分析失败
	JobErrTextAnalysisReadJSONFailed  = 104 // 文本分析JSON读取失败
	JobErrVideoTextAnalysisFailed     = 105 // 文本分析失败
	JobErrVideoAnalysisReadJSONFailed = 106 // 视频分析JSON读取失败
)

// 资源状态值
const (
	ResourceDownloading = 0 // 下载资源
	ResourceExtracting  = 1 // 提取音频
	ResourceCompleted   = 2 // 完成

	ResourceErrDownloadFailed = 100 // 资源下载失败
	ResourceErrExtractFailed  = 101 // 音频提取失败
)

// MessageJSON 用户利用 POST 提交的数据, 用于为任务创建唯一的 ID
type MessageJSON struct {
	DeviceID string   `json:"deviceid"`
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
		err = fmt.Errorf("%s", gin.H{"error": "Wrong Format"})
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
