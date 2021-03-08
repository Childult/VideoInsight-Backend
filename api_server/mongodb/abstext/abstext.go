package abstext

// AbsText 文本摘要
type AbsText struct {
	Hash     string   `bson:"hash"                json:"hash"`                // 文本摘要的哈希值
	URL      string   `bson:"url"                 json:"url"`                 // 对应链接地址
	KeyWords []string `bson:"key_words,omitempty" json:"key_words,omitempty"` // 关键字
	Text     string   `bson:"text"                json:"text"`                // 语音识别
	Abstract string   `bson:"abstract"            json:"abstract"`            // 摘要
}

// GetKeyTag 返回主键标签
func (media AbsText) GetKeyTag() string {
	return "hash"
}

// GetKeyValue 返回主键值
func (media AbsText) GetKeyValue() string {
	return media.Hash
}

// GetCollName 返回数据库名称
func (media AbsText) GetCollName() string {
	return "abstract_text"
}
