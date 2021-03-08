package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// DeleteOne 通过数据主键删除对应数据
func DeleteOne(document Key) (err error) {
	// 检查数据是否存在
	exists := HaveExisted(document)
	if exists == false {
		return fmt.Errorf("Not Found <%s>", document)
	}

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

	// 删除
	_, err = coll.DeleteOne(ctx, bson.M{document.GetKeyTag(): document.GetKeyValue()})
	if err != nil {
		return fmt.Errorf("Failed to delete <%s>", document)
	}
	return
}

// DeleteOneByfilter 通过设定的过滤条件, 删除搜索到的第一条数据
func DeleteOneByfilter(collName string, filter interface{}) (err error) {
	// 设置连接时间阈值, 这段时间内连接失败会重新尝试
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 初始化数据库
	dba := InitDB()
	dba.Connect()
	defer dba.Disconnect()

	// 获取 media collection 的句柄
	coll := dba.GetCollection(collName)

	// 删除
	_, err = coll.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("Failed to delete")
	}
	return
}

// DeleteManyByfilter 通过设定的过滤条件, 删除搜索到的所有数据
func DeleteManyByfilter(collName string, filter interface{}) (err error) {
	// 设置连接时间阈值, 这段时间内连接失败会重新尝试
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 初始化数据库
	dba := InitDB()
	dba.Connect()
	defer dba.Disconnect()

	// 获取 media collection 的句柄
	coll := dba.GetCollection(collName)

	// 删除
	_, err = coll.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("Failed to delete")
	}
	return
}
