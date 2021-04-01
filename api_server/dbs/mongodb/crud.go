package mongodb

import (
	"context"
	"swc/util"

	"go.mongodb.org/mongo-driver/bson"
)

// Exists 根据主键判断数据是否存在
func Exists(document mongoData) bool {
	coll := Get(util.MongoDB, document.Coll())                                       // 获取 collection 的句柄
	result := coll.FindOne(context.TODO(), bson.M{document.Tag(): document.Value()}) // 查找
	return result.Err() == nil
}

// InsertOne 通过主键插入数据
func InsertOne(document mongoData) (err error) {
	coll := Get(util.MongoDB, document.Coll())        // 获取 collection 的句柄
	_, err = coll.InsertOne(context.TODO(), document) // 插入
	return
}

// DeleteOne 删除数据
func DeleteOne(document mongoData) (err error) {
	coll := Get(util.MongoDB, document.Coll()) // 获取 collection 的句柄
	_, err = coll.DeleteOne(context.TODO(), bson.M{document.Tag(): document.Value()})
	return
}

// UpdataOne 取代主键相同的数据
func UpdataOne(document mongoData) (err error) {
	coll := Get(util.MongoDB, document.Coll())                                                   // 获取 collection 的句柄
	_, err = coll.ReplaceOne(context.TODO(), bson.M{document.Tag(): document.Value()}, document) // 代替数据
	return
}

// FindOne 返回找到的第一个数据, 以主键查找
func FindOne(document mongoData) error {
	coll := Get(util.MongoDB, document.Coll())                                                     // 获取 collection 的句柄
	return coll.FindOne(context.TODO(), bson.M{document.Tag(): document.Value()}).Decode(document) // 查找, 原地填充, 返回
}
