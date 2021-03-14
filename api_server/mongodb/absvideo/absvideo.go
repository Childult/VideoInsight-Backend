package absvideo

import "swc/mongodb"

// AbsVideo 视频摘要
type AbsVideo struct {
	URL      string   `bson:"url"                 json:"url"`      // 对应链接地址
	Abstract []string `bson:"abstract"            json:"abstract"` // 视频摘要地址, 每一项是一张图片
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

// HaveAbsVideoExisted 根据判断是否已经存在
func HaveAbsVideoExisted(url string) bool {
	av := AbsVideo{URL: url}
	return mongodb.HaveExisted(av)
}
