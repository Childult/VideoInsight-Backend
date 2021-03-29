package resource

import (
	"context"
	"fmt"
	"swc/dbs/mongodb"
	"swc/logger"

	"go.mongodb.org/mongo-driver/bson"
)

// ExistInMongodb 通过数据的主键, 查看数据是否存在
func (r *Resource) ExistInMongodb() (b bool) {
	// 获取 media collection 的句柄
	coll := mongodb.Get(Database, Collection)
	result := coll.FindOne(context.TODO(), bson.M{r.GetKeyTag(): r.GetKeyValue()})
	return result.Err() == nil
}

// Dump 将数据持久化到 mongodb 中
func (r *Resource) Dump() (err error) {
	// 检查数据是否存在
	if r.ExistInMongodb() {
		// 存在则更新
		coll := mongodb.Get(Database, Collection)                                            // 获取 media collection 的句柄
		_, err = coll.ReplaceOne(context.TODO(), bson.M{r.GetKeyTag(): r.GetKeyValue()}, *r) // 更新
	} else {
		// 不存在则插入
		coll := mongodb.Get(Database, Collection)   // 获取 media collection 的句柄
		_, err = coll.InsertOne(context.TODO(), *r) // 插入
	}

	if err != nil {
		logger.Error.Println(err.Error())
		err = fmt.Errorf("数据<%v>插入失败", *r)
		return err
	}
	return err
}

// Load 从 mongodb 中加载数据
func (r *Resource) Load() (err error) {
	// 检查数据是否存在
	if r.ExistInMongodb() {
		coll := mongodb.Get(Database, Collection) // 获取 collection 的句柄

		// 加载数据
		resource := Resource{}
		err = coll.FindOne(context.TODO(), bson.M{r.GetKeyTag(): r.GetKeyValue()}).Decode(&resource)
		*r = resource
		if err != nil {
			logger.Error.Println(err.Error())
			err = fmt.Errorf("未知错误<%v>", err)
		}
		return err
	} else {
		err = fmt.Errorf("未找到<%v>", r)
		return err
	}
}

// Delete 从 mongodb 中删除数据
func (r *Resource) Delete() (err error) {
	// 检查数据是否存在
	if r.ExistInMongodb() {
		coll := mongodb.Get(Database, Collection) // 获取 collection 的句柄
		// 删除
		_, err = coll.DeleteOne(context.TODO(), bson.M{r.GetKeyTag(): r.GetKeyValue()})
		if err != nil {
			logger.Error.Println(err.Error())
			err = fmt.Errorf("删除<%v>失败", r)
		}
	}

	return err
}
