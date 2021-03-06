package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// DeleteOne as indicated by the name
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
	dba := initDB()
	dba.connect()
	defer dba.disconnect()

	// 获取 media collection 的句柄
	collName := document.GetCollName()
	coll := dba.getCollection(collName)

	// 删除
	_, err = coll.DeleteOne(ctx, bson.M{document.GetKeyTag(): document.GetKeyValue()})
	if err != nil {
		return fmt.Errorf("Failed to delete <%s>", document)
	}
	return
}

// DeleteOneByfilter as indicated by the name
func DeleteOneByfilter(collName string, filter interface{}) (err error) {
	// 设置连接时间阈值, 这段时间内连接失败会重新尝试
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 初始化数据库
	dba := initDB()
	dba.connect()
	defer dba.disconnect()

	// 获取 media collection 的句柄

	coll := dba.getCollection(collName)

	// 删除
	_, err = coll.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("Failed to delete")
	}
	return
}

// DeleteManyByfilter as indicated by the name
func DeleteManyByfilter(collName string, filter interface{}) (err error) {
	// 设置连接时间阈值, 这段时间内连接失败会重新尝试
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 初始化数据库
	dba := initDB()
	dba.connect()
	defer dba.disconnect()

	// 获取 media collection 的句柄

	coll := dba.getCollection(collName)

	// 删除
	_, err = coll.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("Failed to delete")
	}
	return
}