package abstext

import (
	"crypto/sha1"
	"fmt"
	"strings"
)

const (
	Collection = "abstract_text"
)

// AbsText 文本摘要
type AbsText struct {
	Hash     string   `bson:"hash"                json:"hash"`                // 文本摘要的哈希值
	URL      string   `bson:"url"                 json:"url"`                 // 对应链接地址
	KeyWords []string `bson:"key_words"           json:"key_words,omitempty"` // 关键字
	Text     string   `bson:"text"                json:"text"`                // 语音识别
	Abstract []string `bson:"abstract"            json:"abstract"`            // 摘要
	Status   int32    `bson:"status"             json:"status"`               // 当前状态
}

// Tag 返回主键标签
func (at *AbsText) Tag() string {
	return "hash"
}

// Value 返回主键值
func (at *AbsText) Value() string {
	return at.Hash
}

// Coll 返回表名
func (at *AbsText) Coll() string {
	return Collection
}

func NewAbsText(url string, keyWords []string) (at *AbsText) {
	at = &AbsText{
		Hash:     getAbsTextHash(url, keyWords),
		URL:      url,
		KeyWords: keyWords,
	}
	return
}

// getAbsTextHash 从 url 和 keywords 中获取哈希
func getAbsTextHash(url string, keyWords []string) string {
	var str [12]byte                            // 固定大小
	hash := sha1.New()                          // 使用 sha1 哈希函数
	textStr := url + strings.Join(keyWords, "") // 根据  url 和关键词生成哈希
	hash.Write([]byte(textStr))                 // 写入哈希函数
	copy(str[:], hash.Sum([]byte(""))[0:12])    // 复制

	return fmt.Sprintf("%x", str) // 转为字符串
}
