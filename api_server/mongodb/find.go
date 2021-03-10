package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// FindOneByfilter 通过过滤条件和集合名, 返回搜索到的第一条数据
func FindOneByfilter(collName string, filter interface{}) (data bson.M, err error) {
	// 设置连接时间阈值, 这段时间内连接失败会重新尝试
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 初始化数据库
	dba := InitDB()
	dba.Connect()
	defer dba.Disconnect()

	// 获取 media collection 的句柄
	coll := dba.GetCollection(collName)

	// 搜索
	err = coll.FindOne(ctx, filter).Decode(&data)
	if err != nil {
		err = fmt.Errorf("Not Found <%s>", filter)
	}
	return data, err
}

// FindOne 通过主键, 寻找并返回数据
func FindOne(document Key) (data bson.M, err error) {
	// 设置连接时间阈值, 这段时间内连接失败会重新尝试
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 初始化数据库
	dba := InitDB()
	dba.Connect()
	defer dba.Disconnect()

	// 获取 media collection 的句柄
	collName := document.GetCollName()
	coll := dba.GetCollection(collName)

	// 搜索
	KeyTag := document.GetKeyTag()
	KeyValue := document.GetKeyValue()
	err = coll.FindOne(ctx, bson.M{KeyTag: KeyValue}).Decode(&data)
	if err != nil {
		err = fmt.Errorf("Not Found <%s>", document)
	}
	return data, err
}
