package job

import (
	"context"
	"fmt"
	"swc/mongodb"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// Job 用户每个请求对应一个任务
type Job struct {
	JobID    string   `bson:"job_id"              json:"job_id"`        // 唯一ID, 是一个 hash 值
	DeviceID string   `bson:"device_id"           json:"device_id"`     // 用户设备ID
	URL      string   `bson:"url"                 json:"url"`           // 目标地址
	KeyWords []string `bson:"key_words"           json:"key_words"`     // 用户创建的关键字
	Status   int32    `bson:"status"              json:"status"`        // 当前任务状态
	AbsText  string   `bson:"abstract_text"       json:"abstract_text"` // 文本摘要在数据库中的哈希值, 可以复用
}

// GetKeyTag 返回主键标签
func (j Job) GetKeyTag() string {
	return "job_id"
}

// GetKeyValue 返回主键值
func (j Job) GetKeyValue() string {
	return j.JobID
}

// GetCollName 返回数据库名称
func (j Job) GetCollName() string {
	return "job"
}

// SetStatus 设置状态并更新
func (j *Job) SetStatus(i int32) {
	j.Status = i
	mongodb.Update(j)
}

// SetAbsText 设置文本摘要哈希地址并更新
func (j *Job) SetAbsText(hash string) {
	j.AbsText = hash
	mongodb.Update(j)
}

// GetByKey 通过 JobID 返回数据内容
func GetByKey(JobID string) (j Job, err error) {
	// 设置连接时间阈值, 这段时间内连接失败会重新尝试
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 初始化数据库
	dba := mongodb.InitDB()
	dba.Connect()
	defer dba.Disconnect()

	// 获取 media collection 的句柄
	collName := j.GetCollName()
	coll := dba.GetCollection(collName)

	// 搜索
	key := j.GetKeyTag()
	value := JobID
	filter := bson.M{key: value}
	err = coll.FindOne(ctx, filter).Decode(&j)
	if err != nil {
		err = fmt.Errorf("Not Found <%s>", filter)
	}
	return
}
