package mongodb

import (
	"context"
	"swc/logger"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// DBAccessor 数据库访问对象
type DBAccessor struct {
	client *mongo.Client
}

// Key 实现该接口则能对集合进行访问
type Key interface {
	GetKeyTag() string
	GetKeyValue() string
	GetCollName() string
}

var (
	// SWCDB 数据库名
	SWCDB = "swcdb"
)

const (
	// MongoDBURI mongodb 地址
	MongoDBURI = "mongodb://192.168.2.80:27018"
)

// InitDB 返回 DBAccessor
func InitDB() (dba *DBAccessor) {
	dba = &DBAccessor{}
	return
}

// Connect 连接到 mongodb
func (dba *DBAccessor) Connect() (err error) {
	// 创建 mongodb client
	dba.client, err = mongo.NewClient(options.Client().ApplyURI(MongoDBURI))
	if err != nil {
		logger.Error.Println(err.Error())
	}

	// 设置连接时间阈值, 这段时间内连接失败会重新尝试
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// 利用 client 连接 mongodb, 超时会产生错误
	err = dba.client.Connect(ctx)
	if err != nil {
		logger.Error.Println(err.Error())
	}

	// ping 测试
	err = dba.client.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.Error.Println(err.Error())
	}

	return nil
}

// Disconnect 断开连接, mongo-driver 中实现了连接池, 所以每次用完后即时断开即可
func (dba *DBAccessor) Disconnect() (err error) {
	// 设置连接时间阈值, 这段时间内连接失败会重新尝试
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = dba.client.Disconnect(ctx)
	if err != nil {
		logger.Error.Println(err.Error())
	}

	return nil
}

// GetCollection 通过集合名称获取访问该集合的句柄
func (dba *DBAccessor) GetCollection(name string) (coll *mongo.Collection) {

	// 连接数据库
	db := dba.client.Database(SWCDB)

	// 获取 collection 的句柄
	coll = db.Collection(name)

	return
}

// ShowAllDatabaseNames 显示所有数据库
func (dba *DBAccessor) ShowAllDatabaseNames() (databases []string) {
	// 设置连接时间阈值, 这段时间内连接失败会重新尝试
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	databases, err := dba.client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		logger.Error.Println(err.Error())
	}
	return
}

// HaveExisted 通过数据的主键, 查看数据是否存在
func HaveExisted(document Key) (b bool) {
	// 设置连接时间阈值, 这段时间内连接失败会重新尝试
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 初始化数据库
	dba := InitDB()
	dba.Connect()
	defer dba.Disconnect()

	// 获取 collName collection 的句柄
	collName := document.GetCollName()
	coll := dba.GetCollection(collName)
	result := coll.FindOne(ctx, bson.M{document.GetKeyTag(): document.GetKeyValue()})
	if result.Err() == nil {
		return true
	}
	return false
}
