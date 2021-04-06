package mongodb

import (
	"context"
	"fmt"
	"swc/logger"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	poolSize    = 200              // 维持的最大连接数
	idleTimeout = 30 * time.Second // 每个连接维持的最长时间
)

type mongodbPool struct {
	*mongo.Client
}

var (
	pool mongodbPool // mongodb 连接池
)

// 初始化 mongodb 信息, 可以从别的地方进行配置
func InitMongodb(addr, user, password string) {
	var err error
	var uri string
	opts := options.Client()
	if user != "" && password != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s", user, password, addr)
	} else {
		uri = fmt.Sprintf("mongodb://%s", addr)
	}

	opts.ApplyURI(uri)                                     // mongodb 连接地址
	opts.SetMaxPoolSize(poolSize)                          // 最大连接数
	opts.SetMaxConnIdleTime(idleTimeout)                   // 每个连接持续最长时间
	pool.Client, err = mongo.Connect(context.TODO(), opts) // 创建全局连接
	if err != nil {
		logger.Error.Fatal("mongodb connect failed")
	}
}

// Get 根据 dbName 和 collName, 返回 *mongo.Collection
// dbName: 需要访问的 mongodb 数据库名称
// collName: 需要访问的 mongodb 集合名词
// coll: 对应集合的句柄
func Get(dbName, collName string) (coll *mongo.Collection) {
	return pool.Database(dbName).Collection(collName)
}

// mongoData 定义 mongodb 数据的接口, 实现该接口可以在 mongodb 里增删改查
// 需要知道主键名, 主键值, 表名
type mongoData interface {
	Tag() string
	Value() string
	Coll() string
}
