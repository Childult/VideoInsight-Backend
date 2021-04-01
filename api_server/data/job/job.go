package job

import (
	"crypto/sha1"
	"fmt"
	"strings"
	"swc/util"
)

const (
	Database   = "swcdb"
	Collection = "job"
)

// Job 用户每个请求对应一个任务
type Job struct {
	JobID    string   `bson:"job_id"              json:"job_id"`        // 唯一ID, 是一个 hash 值
	DeviceID string   `bson:"device_id"           json:"device_id"`     // 用户设备ID
	URL      string   `bson:"url"                 json:"url"`           // 目标地址
	KeyWords []string `bson:"key_words"           json:"key_words"`     // 用户创建的关键字
	Status   int32    `bson:"status"              json:"status"`        // 当前任务状态
	AbsText  string   `bson:"abstract_text"       json:"abstract_text"` // 文本摘要在数据库中的哈希值, 可以复用
}

// Tag 返回主键名
func (j *Job) Tag() string {
	return "job_id"
}

// Value 返回主键值
func (j *Job) Value() string {
	return j.JobID
}

// Coll 返回表名
func (j *Job) Coll() string {
	return Collection
}

func NewJob(dID, url string, keyWords []string) (j *Job) {
	j = &Job{
		JobID:    getJobID(dID, url, keyWords),
		DeviceID: dID,
		URL:      url,
		KeyWords: keyWords,
		Status:   util.JobStart,
	}
	if len(keyWords) == 0 {
		j.KeyWords = nil
	}
	return
}

// getJobID 从 url 和 keywords 中获取哈希
func getJobID(dID, url string, keyWords []string) string {
	var str [12]byte                                  // 固定大小
	hash := sha1.New()                                // 使用 sha1 哈希函数
	textStr := dID + url + strings.Join(keyWords, "") // 根据 deviceID, url 和关键词生成哈希
	hash.Write([]byte(textStr))                       // 写入哈希函数
	copy(str[:], hash.Sum([]byte(""))[0:12])          // 复制

	return fmt.Sprintf("%x", str) // 转为字符串
}
