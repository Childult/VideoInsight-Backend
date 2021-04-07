package resource

const (
	Collection = "resource"
)

// Resource 目标资源
type Resource struct {
	URL       string `bson:"url"                json:"url"`        // 逐渐
	Location  string `bson:"location"           json:"location"`   // 存储路径
	VideoPath string `bson:"video_path"         json:"video_path"` // 视频文件名
	AudioPath string `bson:"audio_path"         json:"audio_path"` // 音频文件名
	Title     string `bson:"title"              json:"title"`      // 标题
	Status    int32  `bson:"status"             json:"status"`     // 当前状态
}

// Tag 返回主键标签
func (r *Resource) Tag() string {
	return "url"
}

// Value 返回主键值
func (r *Resource) Value() string {
	return r.URL
}

// Coll 返回表名
func (r *Resource) Coll() string {
	return Collection
}
