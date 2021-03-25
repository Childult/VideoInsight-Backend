package resource

const (
	Database   = "swcdb"
	Collection = "resource"
)

// Resource 目标资源
type Resource struct {
	URL       string `bson:"url"                json:"url"`           // 逐渐
	Status    int32  `bson:"status"             json:"status"`        // 当前状态
	Location  string `bson:"location"           json:"location"`      // 存储路径
	VideoPath string `bson:"video_path"         json:"video_path"`    // 视频文件名
	AudioPath string `bson:"audio_path"         json:"audio_path"`    // 音频文件名
	AbsText   string `bson:"abstract_text"      json:"abstract_text"` // 无关键词对应的文本摘要哈希
}

// GetKeyTag 返回主键标签
func (r *Resource) GetKeyTag() string {
	return "url"
}

// GetKeyValue 返回主键值
func (r *Resource) GetKeyValue() string {
	return r.URL
}
