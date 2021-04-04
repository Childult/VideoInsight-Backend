package absvideo

const (
	Collection = "abstract_video"
)

// AbsVideo 视频摘要
type AbsVideo struct {
	URL      string   `bson:"url"                 json:"url"`      // 对应链接地址
	Abstract []string `bson:"abstract"            json:"abstract"` // 视频摘要地址, 每一项是一张图片
	Status   int32    `bson:"status"             json:"status"`    // 当前状态
}

// Tag 返回主键标签
func (av *AbsVideo) Tag() string {
	return "url"
}

// Value 返回主键值
func (av *AbsVideo) Value() string {
	return av.URL
}

// Coll 返回表名
func (av *AbsVideo) Coll() string {
	return Collection
}
