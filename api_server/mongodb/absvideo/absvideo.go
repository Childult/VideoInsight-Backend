package absvideo

// AbsVideo 视频摘要
type AbsVideo struct {
	Hash     string   `bson:"hash"                json:"hash"`                // 文本摘要的哈希值
	URL      string   `bson:"url"                 json:"url"`                 // 对应链接地址
	KeyWords []string `bson:"key_words,omitempty" json:"key_words,omitempty"` // 关键字
	Text     string   `bson:"text"                json:"text"`                // 语音识别
	Abstract string   `bson:"abstract"            json:"abstract"`            // 摘要
}

// GetKeyTag 返回主键标签
func (media AbsVideo) GetKeyTag() string {
	return "hash"
}

// GetKeyValue 返回主键值
func (media AbsVideo) GetKeyValue() string {
	return media.Hash
}

// GetCollName 返回数据库名称
func (media AbsVideo) GetCollName() string {
	return "abstract_text"
}
