package job

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
func (j *Job) Dump() (err error) {
	// 检查数据是否存在
	if haveExisted(j) {
		// 存在则更新
		coll := mongodb.Get(Database, Collection)                                            // 获取 media collection 的句柄
		_, err = coll.ReplaceOne(context.TODO(), bson.M{j.GetKeyTag(): j.GetKeyValue()}, *j) // 更新
	} else {
		// 不存在则插入
		coll := mongodb.Get(Database, Collection)   // 获取 media collection 的句柄
		_, err = coll.InsertOne(context.TODO(), *j) // 插入
	}

	if err != nil {
		logger.Error.Println(err.Error())
		err = fmt.Errorf("数据<%v>插入失败", *j)
		return err
	}
	return err
}

// Load 从 mongodb 中加载数据
func (j *Job) Load() (err error) {
	coll := mongodb.Get(Database, Collection) // 获取 collection 的句柄

	// 加载数据
	job := Job{}
	err = coll.FindOne(context.TODO(), bson.M{j.GetKeyTag(): j.GetKeyValue()}).Decode(&job)
	*j = job
	if err != nil {
		logger.Error.Println(err.Error())
		err = fmt.Errorf("未找到<%v>", j)
	}
	return err
}

// Delete 从 mongodb 中删除数据
func (j *Job) Delete() (err error) {
	// 检查数据是否存在
	if haveExisted(j) {
		coll := mongodb.Get(Database, Collection) // 获取 collection 的句柄
		// 删除
		_, err = coll.DeleteOne(context.TODO(), bson.M{j.GetKeyTag(): j.GetKeyValue()})
		if err != nil {
			logger.Error.Println(err.Error())
			err = fmt.Errorf("删除<%v>失败", j)
		}
	}

	return err
}
