package abstext

import (
	"crypto/sha1"
	"fmt"
	"strings"
	"swc/mongodb"
)

// AbsText 文本摘要
type AbsText struct {
	Hash     string   `bson:"hash"                json:"hash"`                // 文本摘要的哈希值
	URL      string   `bson:"url"                 json:"url"`                 // 对应链接地址
	KeyWords []string `bson:"key_words,omitempty" json:"key_words,omitempty"` // 关键字
	Text     string   `bson:"text"                json:"text"`                // 语音识别
	Abstract string   `bson:"abstract"            json:"abstract"`            // 摘要
}

// GetKeyTag 返回主键标签
func (at AbsText) GetKeyTag() string {
	return "hash"
}

// GetKeyValue 返回主键值
func (at AbsText) GetKeyValue() string {
	return at.Hash
}

// GetCollName 返回数据库名称
func (at AbsText) GetCollName() string {
	return "abstract_text"
}

// NewAbsText 构建结构体
func NewAbsText(url, text, abstract string, keyWords []string) (at AbsText) {
	at.URL = url
	at.Text = text
	at.Abstract = abstract
	at.KeyWords = keyWords
	at.Hash = at.getAbsHash()
	return
}

// 从 url 和 keywords 中获取哈希
func (at *AbsText) getAbsHash() string {
	var str [12]byte                                  // 固定大小
	hash := sha1.New()                                // 使用 sha1 哈希函数
	textStr := at.URL + strings.Join(at.KeyWords, "") // 根据 url 和关键字生成哈市
	hash.Write([]byte(textStr))                       // 写入哈希函数
	copy(str[:], hash.Sum([]byte(""))[0:12])          // 复制

	return fmt.Sprintf("%x", str) // 转为字符串
}

// HaveAbsTextExisted 根据判断是否已经存在
func HaveAbsTextExisted(url string, keyWords []string) bool {
	at := NewAbsText(url, "", "", keyWords)
	return mongodb.HaveExisted(at)
}
