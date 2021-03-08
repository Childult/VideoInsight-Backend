package resource

import (
	"context"
	"fmt"
	"swc/mongodb"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// Resource include media and audio
type Resource struct {
	URL       string `bson:"url"                json:"url"`
	Status    int32  `bson:"status"             json:"status"`
	Location  string `bson:"location"           json:"location"`
	VideoPath string `bson:"video_path"         json:"video_path"`
	AudioPath string `bson:"audio_path"         json:"audio_path"`
	AbsText   string `bson:"abstract_text"      json:"abstract_text"`
	AbsVideo  string `bson:"abstract_video"     json:"abstract_video"`
}

// GetKeyTag implement the interface Key
func (r Resource) GetKeyTag() string {
	return "url"
}

// GetKeyValue implement the interface Key
func (r Resource) GetKeyValue() string {
	return r.URL
}

// GetCollName implement the interface Key
func (r Resource) GetCollName() string {
	return "resource"
}

// GetByKey return one resource
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

// Refresh read the database, refresh variables in place
func (s *Resource) Refresh() (err error) {
	*s, err = GetByKey(s.URL)
	return
}

// SetStatus implement the interface Key
func (s *Resource) SetStatus(status int32) {
	s.Status = status
	mongodb.Update(s)
}

// SetAbsText implement the interface Key
func (s *Resource) SetAbsText(key string) {
	s.AbsText = key
	mongodb.Update(s)
}

// SetAbsVideo implement the interface Key
func (s *Resource) SetAbsVideo(key string) {
	s.AbsVideo = key
	mongodb.Update(s)
}
