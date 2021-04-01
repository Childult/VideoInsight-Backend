package mongodb

import (
	"context"
	"fmt"
	"swc/util"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func init() {
	InitMongodb("192.168.2.80:27018", "", "")
	util.MongoDB = "test"
}

func TestShowDatabases(t *testing.T) {
	// 查看所有数据库, 用于简单测试是否连接上数据库
	client := pool.Client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer client.Disconnect(ctx)
	err := client.Ping(ctx, readpref.Primary())
	if err != nil {
		t.Error(err)
	}
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(databases)
}

func TestShowCollections(t *testing.T) {
	// 查看所有数据表
	client := pool.Client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer client.Disconnect(ctx)
	err := client.Ping(ctx, readpref.Primary())
	if err != nil {
		t.Error(err)
	}
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		t.Error(err)
	}
	fmt.Println("数据库:", databases)
	for _, dbName := range databases {
		db := client.Database(dbName)
		colls, err := db.ListCollectionNames(ctx, bson.M{})
		if err != nil {
			t.Error(err)
		}
		fmt.Println(dbName, colls)
	}
}

func TestPollConnect(t *testing.T) {
	// mongodb 连接池测试, 测试连接池设置是否正确
	// 测试运行后可以在 mongodb 终端中输入 db.serverStatus().connections, 查看连接数
	coll := Get("test", "test")
	sw := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		sw.Add(1)
		go func(sw *sync.WaitGroup, coll *mongo.Collection) {
			ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
			defer cancel()
			coll.Find(ctx, bson.M{})
			time.Sleep(time.Second * 10)
			sw.Done()
		}(&sw, coll)
	}
	sw.Wait()
}

// fakeData 测试数据
type fakeData struct {
	Addr string `bson:"addr"`
	Name string `bson:"name"`
}

func (r *fakeData) Tag() string {
	return "addr"
}
func (r *fakeData) Value() string {
	return r.Addr
}
func (r *fakeData) Coll() string {
	return "fakeData"
}

func TestFindAll(t *testing.T) {
	// 要查看的数据类型
	var d *fakeData = new(fakeData)

	// 查看所有数据
	client := pool.Client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer client.Disconnect(ctx)
	err := client.Ping(ctx, readpref.Primary())
	if err != nil {
		t.Error(err)
	}
	coll := client.Database(util.MongoDB).Collection(d.Coll())
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		t.Error(err)
	}
	for cursor.Next(ctx) {
		if err = cursor.Decode(d); err != nil {
			t.Error(err)
		}
		fmt.Printf("%+v\n", d)
	}
}

func TestDeleteAll(t *testing.T) {
	// 删除所有测试数据
	var d *fakeData = new(fakeData) // 要删除的数据类型
	// 查看所有数据
	client := pool.Client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer client.Disconnect(ctx)
	err := client.Ping(ctx, readpref.Primary())
	if err != nil {
		t.Error(err)
	}
	coll := client.Database(util.MongoDB).Collection(d.Coll())
	_, err = coll.DeleteMany(ctx, bson.M{})
	if err != nil {
		t.Error(err)
	}
}

func TestCRUD(t *testing.T) {
	var expectErr, actualErr error
	var expectBool, actualBool bool
	var d *fakeData = new(fakeData)

	// 数据一开始不存在
	expectBool = false
	actualBool = Exists(d)
	assert.Equal(t, expectBool, actualBool)

	// 不存在时查找
	expectErr = fmt.Errorf("mongo: no documents in result")
	actualErr = FindOne(d)
	assert.Equal(t, expectErr, actualErr)

	// 插入数据, 空值也是值
	expectErr = nil
	actualErr = InsertOne(d)
	assert.Equal(t, expectErr, actualErr)

	// 插入后存在
	expectBool = true
	actualBool = Exists(d)
	assert.Equal(t, expectBool, actualBool)

	// 根据主键更新数据
	d.Name = "鸡排卷饼"
	expectErr = nil
	actualErr = UpdataOne(d)
	assert.Equal(t, expectErr, actualErr)

	// 存在时查找
	var dd *fakeData = new(fakeData)
	expectErr = nil
	actualErr = FindOne(dd)
	assert.Equal(t, expectErr, actualErr)
	assert.Equal(t, d, dd)

	// 删除数据
	expectErr = nil
	actualErr = DeleteOne(dd)
	assert.Equal(t, expectErr, actualErr)
}
