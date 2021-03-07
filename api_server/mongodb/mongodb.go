package mongodb

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// DBAccessor is a mongodb accessor
type DBAccessor struct {
	client *mongo.Client
}

// Key return a key
type Key interface {
	GetKeyTag() string
	GetKeyValue() string
	GetCollName() string
}

var (
	// SWCDB swagger codegen db
	SWCDB = "swcdb"
)

const (
	// MongoDBURI mongodb 地址
	MongoDBURI = "mongodb://localhost:27017"
)

// initDB return one *DBAccessor
func initDB() (dba *DBAccessor) {
	dba = &DBAccessor{}
	return
}

// connect to mongodb
func (dba *DBAccessor) connect() (err error) {
	// 创建 mongodb client
	dba.client, err = mongo.NewClient(options.Client().ApplyURI(MongoDBURI))
	if err != nil {
		log.Fatal(err)
	}

	// 设置连接时间阈值, 这段时间内连接失败会重新尝试
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// 利用 client 连接 mongodb, 超时会产生错误
	err = dba.client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// ping 测试
	err = dba.client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// disconnect to mongodb
func (dba *DBAccessor) disconnect() (err error) {
	// 设置连接时间阈值, 这段时间内连接失败会重新尝试
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = dba.client.Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// getCollection return a collection according to name
func (dba *DBAccessor) getCollection(name string) (coll *mongo.Collection) {

	// 连接数据库
	db := dba.client.Database(SWCDB)

	// 获取 collection 的句柄
	coll = db.Collection(name)

	return
}

// ShowAllDatabaseNames show all dbs
func (dba *DBAccessor) ShowAllDatabaseNames() (databases []string) {
	// 设置连接时间阈值, 这段时间内连接失败会重新尝试
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	databases, err := dba.client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	return
}

// HaveExisted as indicated by the name
func HaveExisted(document Key) (b bool) {
	// 设置连接时间阈值, 这段时间内连接失败会重新尝试
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 初始化数据库
	dba := initDB()
	dba.connect()
	defer dba.disconnect()

	// 获取 collName collection 的句柄
	collName := document.GetCollName()
	coll := dba.getCollection(collName)
	result := coll.FindOne(ctx, bson.M{document.GetKeyTag(): document.GetKeyValue()})
	if result.Err() == nil {
		return true
	}
	return false
}
