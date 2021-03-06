package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// FindOneByfilter as indicated by the name
func FindOneByfilter(collName string, filter interface{}) (data bson.M, err error) {
	// 设置连接时间阈值, 这段时间内连接失败会重新尝试
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 初始化数据库
	dba := initDB()
	dba.connect()
	defer dba.disconnect()

	// 获取 media collection 的句柄
	coll := dba.getCollection(collName)

	// 搜索
	err = coll.FindOne(ctx, filter).Decode(&data)
	if err != nil {
		err = fmt.Errorf("Not Found <%s>", filter)
	}
	return data, err
}
