package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// Update replace one
func Update(document Key) (err error) {
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

	// 替换
	KeyTag := document.GetKeyTag()
	KeyValue := document.GetKeyValue()
	_, err = coll.ReplaceOne(ctx, bson.M{KeyTag: KeyValue}, document)
	if err != nil {
		return fmt.Errorf("Failed to update <%s>", document)
	}
	return err
}

// UpdateOne replace one
func UpdateOne(rawData, newData Key) (err error) {
	// 检查数据是否属于同一张表
	collName := rawData.GetCollName()
	newCollName := rawData.GetCollName()
	if collName != newCollName {
		return fmt.Errorf("The collection <%s> does not match with <%s>", collName, newCollName)
	}

	// 检查数据是否存在
	exists := HaveExisted(rawData)
	if exists == false {
		return fmt.Errorf("Not Found <%s>", rawData)
	}

	// 检查数据主键是否相同
	KeyTag := rawData.GetKeyTag()
	KeyValue := rawData.GetKeyValue()
	newKeyTag := newData.GetKeyTag()
	newKeyValue := newData.GetKeyValue()
	if KeyTag != newKeyTag || KeyValue != newKeyValue {
		return fmt.Errorf("The primary key <%s> and <%s> do not match", rawData, newData)
	}

	// 设置连接时间阈值, 这段时间内连接失败会重新尝试
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 初始化数据库
	dba := initDB()
	dba.connect()
	defer dba.disconnect()

	// 获取 media collection 的句柄
	coll := dba.getCollection(collName)

	// 替换
	_, err = coll.ReplaceOne(ctx, bson.M{KeyTag: KeyValue}, newData)
	if err != nil {
		return fmt.Errorf("Failed to update from <%s> to <%s>", rawData, newData)
	}
	return err
}

// ReplaceOneByFilter replace one
func ReplaceOneByFilter(collName string, filter interface{}, update interface{}) (err error) {
	// 设置连接时间阈值, 这段时间内连接失败会重新尝试
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 初始化数据库
	dba := initDB()
	dba.connect()
	defer dba.disconnect()

	// 获取 media collection 的句柄
	coll := dba.getCollection(collName)

	// 替换
	_, err = coll.ReplaceOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("Failed to update from <%s> to <%s>", filter, update)
	}
	return err
}

// UpdateOneByFilter update one
func UpdateOneByFilter(collName string, filter interface{}, update interface{}) (err error) {
	// 设置连接时间阈值, 这段时间内连接失败会重新尝试
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 初始化数据库
	dba := initDB()
	dba.connect()
	defer dba.disconnect()

	// 获取 media collection 的句柄
	coll := dba.getCollection(collName)

	// 替换
	_, err = coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("Failed to update from <%s> to <%s>", filter, update)
	}
	return err
}
