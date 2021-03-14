package resource

import (
	"context"
	"fmt"
	"swc/mongodb"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
func (r Resource) GetKeyTag() string {
	return "url"
}

// GetKeyValue 返回主键值
func (r Resource) GetKeyValue() string {
	return r.URL
}

// GetCollName 返回数据库名称
func (r Resource) GetCollName() string {
	return "resource"
}

// GetByKey 通过 url 返回数据内容
func GetByKey(url string) (s Resource, err error) {
	// 设置连接时间阈值, 这段时间内连接失败会重新尝试
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 初始化数据库
	dba := mongodb.InitDB()
	dba.Connect()
	defer dba.Disconnect()

	// 获取 media collection 的句柄
	collName := s.GetCollName()
	coll := dba.GetCollection(collName)

	// 搜索
	key := s.GetKeyTag()
	value := url
	filter := bson.M{key: value}
	err = coll.FindOne(ctx, filter).Decode(&s)
	if err != nil {
		err = fmt.Errorf("Not Found <%s>", filter)
	}
	return
}

// Refresh 原地更新数据
func (r *Resource) Refresh() (err error) {
	*r, err = GetByKey(r.URL)
	return
}

// SetStatus 设置状态并更新
func (r *Resource) SetStatus(status int32) {
	r.Status = status
	mongodb.Update(r)
}

// SetAbsText 设置文本摘要哈希地址并更新
func (r *Resource) SetAbsText(hash string) {
	r.AbsText = hash
	mongodb.Update(r)
}
