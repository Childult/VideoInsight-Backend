package mongodb

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func init() {
	InitMongodb("172.17.0.3:27017", "", "")
}

func TestShowDatabases(t *testing.T) {
	// 查看所有数据库, 用于简单测试是否连接上数据库
	client := pool.Client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer client.Disconnect(ctx)
	err := client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)
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
