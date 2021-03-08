package mongodb

import (
	"context"
	"fmt"
	"time"
)

// InsertOne as indicated by the name
func InsertOne(document Key) (err error) {
	// 检查数据是否存在
	exists := HaveExisted(document)
	if exists {
		return fmt.Errorf("The document <%s> already exists", document)
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

	// 插入
	_, err = coll.InsertOne(ctx, document)
	if err != nil {
		return fmt.Errorf("The document <%s> insertion failed ", document)
	}
	return err
}
