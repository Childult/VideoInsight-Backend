package absvideo

import (
	"context"
	"fmt"
	"swc/dbs"
	"swc/dbs/mongodb"
	"swc/logger"

	"go.mongodb.org/mongo-driver/bson"
)

// haveExisted 通过数据的主键, 查看数据是否存在
func haveExisted(data dbs.PrimaryKey) (b bool) {
	// 获取 media collection 的句柄
	coll := mongodb.Get(Database, Collection)
	result := coll.FindOne(context.TODO(), bson.M{data.GetKeyTag(): data.GetKeyValue()})
	return result.Err() == nil
}

// Dump 将数据持久化到 mongodb 中
func (av *AbsVideo) Dump() (err error) {
	// 检查数据是否存在
	if haveExisted(av) {
		// 存在则更新
		coll := mongodb.Get(Database, Collection)                                               // 获取 media collection 的句柄
		_, err = coll.ReplaceOne(context.TODO(), bson.M{av.GetKeyTag(): av.GetKeyValue()}, *av) // 更新
	} else {
		// 不存在则插入
		coll := mongodb.Get(Database, Collection)    // 获取 media collection 的句柄
		_, err = coll.InsertOne(context.TODO(), *av) // 插入
	}

	if err != nil {
		logger.Error.Println(err.Error())
		err = fmt.Errorf("数据<%v>插入失败", *av)
		return err
	}
	return err
}

// Load 从 mongodb 中加载数据
func (av *AbsVideo) Load() (err error) {
	coll := mongodb.Get(Database, Collection) // 获取 collection 的句柄

	// 加载数据
	absVideo := AbsVideo{}
	err = coll.FindOne(context.TODO(), bson.M{av.GetKeyTag(): av.GetKeyValue()}).Decode(&absVideo)
	*av = absVideo
	if err != nil {
		logger.Error.Println(err.Error())
		err = fmt.Errorf("未找到<%v>", av)
	}
	return err
}

// Delete 从 mongodb 中删除数据
func (av *AbsVideo) Delete() (err error) {
	// 检查数据是否存在
	if haveExisted(av) {
		coll := mongodb.Get(Database, Collection) // 获取 collection 的句柄
		// 删除
		_, err = coll.DeleteOne(context.TODO(), bson.M{av.GetKeyTag(): av.GetKeyValue()})
		if err != nil {
			logger.Error.Println(err.Error())
			err = fmt.Errorf("删除<%v>失败", av)
		}
	}

	return err
}
