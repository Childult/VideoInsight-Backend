package util

import (
	"crypto/sha1"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	// WorkSpace 动态设置工作路径
	WorkSpace = ""
	// SavePath 文件保存位置
	SavePath = "/home/download"
)

// 任务状态值
const (
	JobStart              = 0 // 创建资源, 写入数据库
	JobDownloadMedia      = 1 // 下载资源
	JobExisted            = 2 // 文件已存在
	JobExtractAudio       = 3 // 提取音频
	JobExtractDone        = 4 // 音频提取成功
	JobAbstractextraction = 5 // 提取摘要
	JobCompleted          = 6 // 完成

	JobErrFailedToFindResource       = 100 // 从数据库中读取时发生错误
	JobErrDownloadFailed             = 101 // 资源下载失败
	JobErrExtractFailed              = 102 // 音频提取失败
	JobErrTextAnalysisFailed         = 103 // 文本分析失败
	JobErrTextAnalysisReadJSONFailed = 104 // 文本分析JSON读取失败
)

// 资源状态值
const (
	ResourceDownloading = 0 // 下载资源
	ResourceExtracting  = 1 // 提取音频
	ResourceCompleted   = 2 // 完成

	ResourceErrDownloadFailed = 100 // 资源下载失败
	ResourceErrExtractFailed  = 101 // 音频提取失败
)

// SetWorkSpace 获取当前路径
func SetWorkSpace() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	WorkSpace, _ = filepath.Split(dir)
}

// MessageJSON contains unique identity and user data
type MessageJSON struct {
	DeviceID string   `json:"deviceid"`
	URL      string   `json:"url"`
	KeyWords []string `json:"keywords,omitempty"`
}

func (json MessageJSON) String() string {
	return json.DeviceID + json.URL + strings.Join(json.KeyWords, "")
}

// GetHash will return a hash
func (json MessageJSON) GetHash() (result [12]byte) {
	hash := sha1.New()
	hash.Write([]byte(json.String()))
	copy(result[:], hash.Sum([]byte(""))[0:12])
	return
}

// GetID will return a id
func (json MessageJSON) GetID() string {
	return fmt.Sprintf("%v", json.GetHash())
}

// GetJSON return a json
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

func removeEmptyString(a []string) []string {
	return deleteKeywords(a, "")
}

func deleteKeywords(a []string, s string) []string {
	j := 0
	for _, val := range a {
		if val == s {
			a[j] = val
			j++
		}
	}
	return a[j:]
}
