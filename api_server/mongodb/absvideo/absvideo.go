package absvideo

// AbsVideo 视频摘要
type AbsVideo struct {
	URL      string   `bson:"url"                 json:"url"`      // 对应链接地址
	Abstract []string `bson:"abstract"            json:"abstract"` // 视频摘要地址
}

// GetKeyTag 返回主键标签
func (av AbsVideo) GetKeyTag() string {
	return "url"
}

// GetKeyValue 返回主键值
func (av AbsVideo) GetKeyValue() string {
	return av.URL
}

// GetCollName 返回数据库名称
func (av AbsVideo) GetCollName() string {
	return "abstract_video"
}
