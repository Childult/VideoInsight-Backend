package task

import (
	"crypto/sha1"
	"fmt"
	"strings"
	"swc/util"
)

const (
	Collection = "task"
)

// Task 用户每个请求对应一个任务
type Task struct {
	TaskID   string   `bson:"task_id"             json:"task_id"`       // 唯一ID, 是一个 hash 值, 也是文本摘的主键
	URL      string   `bson:"url"                 json:"url"`           // 目标地址, 资源和视频摘要的主键
	KeyWords []string `bson:"key_words"           json:"key_words"`     // 用户创建的关键字
	Status   int32    `bson:"status"              json:"status"`        // 当前任务状态
	TextHash string   `bson:"abstract_text"       json:"abstract_text"` // 文本摘要在数据库中的哈希值, 是其主键
}

// Tag 返回主键名
func (t *Task) Tag() string {
	return "task_id"
}

// Value 返回主键值
func (t *Task) Value() string {
	return t.TaskID
}

// Coll 返回表名
func (t *Task) Coll() string {
	return Collection
}

func NewTask(url string, keyWords []string) (j *Task) {
	j = &Task{
		TaskID:   getTaskID(url, nil),
		URL:      url,
		KeyWords: keyWords,
		Status:   util.TaskStart,
	}
	if len(keyWords) == 0 {
		j.KeyWords = nil
	}
	return
}

// getTaskID 从 url 和 keywords 中获取哈希
func getTaskID(url string, keyWords []string) string {
	var str [12]byte                            // 固定大小
	hash := sha1.New()                          // 使用 sha1 哈希函数
	textStr := url + strings.Join(keyWords, "") // 根据 url 和关键词生成哈希
	hash.Write([]byte(textStr))                 // 写入哈希函数
	copy(str[:], hash.Sum([]byte(""))[0:12])    // 复制

	return fmt.Sprintf("%x", str) // 转为字符串
}
