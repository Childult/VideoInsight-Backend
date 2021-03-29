package abstext

import (
	"context"
	"fmt"
	"swc/dbs/mongodb"
	"swc/logger"

	"go.mongodb.org/mongo-driver/bson"
)

// haveExisted 通过数据的主键, 查看数据是否存在
func (at *AbsText) ExistInMongodb() (b bool) {
	// 获取 media collection 的句柄
	coll := mongodb.Get(Database, Collection)
	result := coll.FindOne(context.TODO(), bson.M{at.GetKeyTag(): at.GetKeyValue()})
	return result.Err() == nil
}

// Dump 将数据持久化到 mongodb 中
func (at *AbsText) Dump() (err error) {
	// 检查数据是否存在
	if at.ExistInMongodb() {
		// 存在则更新
		coll := mongodb.Get(Database, Collection)                                               // 获取 media collection 的句柄
		_, err = coll.ReplaceOne(context.TODO(), bson.M{at.GetKeyTag(): at.GetKeyValue()}, *at) // 更新
	} else {
		// 不存在则插入
		coll := mongodb.Get(Database, Collection)    // 获取 media collection 的句柄
		_, err = coll.InsertOne(context.TODO(), *at) // 插入
	}

	if err != nil {
		logger.Error.Println(err.Error())
		err = fmt.Errorf("数据<%v>插入失败", *at)
		return err
	}
	return err
}

// Load 从 mongodb 中加载数据
func (at *AbsText) Load() (err error) {
	// 检查数据是否存在
	if at.ExistInMongodb() {
		coll := mongodb.Get(Database, Collection) // 获取 collection 的句柄

		// 加载数据
		absText := AbsText{}
		err = coll.FindOne(context.TODO(), bson.M{at.GetKeyTag(): at.GetKeyValue()}).Decode(&absText)
		*at = absText
		if err != nil {
			logger.Error.Println(err.Error())
			err = fmt.Errorf("未知<%v>", err)
		}
		return err
	} else {
		err = fmt.Errorf("未找到<%v>", at)
		return err
	}
}

// Delete 从 mongodb 中删除数据
func (at *AbsText) Delete() (err error) {
	// 检查数据是否存在
	if at.ExistInMongodb() {
		coll := mongodb.Get(Database, Collection) // 获取 collection 的句柄
		// 删除
		_, err = coll.DeleteOne(context.TODO(), bson.M{at.GetKeyTag(): at.GetKeyValue()})
		if err != nil {
			logger.Error.Println(err.Error())
			err = fmt.Errorf("删除<%v>失败", at)
		}
	}

	return err
}
