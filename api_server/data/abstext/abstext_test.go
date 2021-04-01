package abstext

import (
	"context"
	"fmt"
	"swc/dbs/mongodb"
	"swc/dbs/redis"
	"swc/util"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	redis.InitRedis("172.17.0.4:6379", "")
	mongodb.InitMongodb("192.168.2.80:27018", "", "")
	util.MongoDB = "test"
}

var (
	d  *AbsText = new(AbsText)
	dd *AbsText = new(AbsText)
)

func TestFindAll(t *testing.T) {
	// 查看所有数据
	coll := mongodb.Get(util.MongoDB, d.Coll())
	cursor, err := coll.Find(context.TODO(), bson.M{})
	if err != nil {
		t.Error(err)
	}
	for cursor.Next(context.TODO()) {
		if err = cursor.Decode(d); err != nil {
			t.Error(err)
		}
		fmt.Printf("%+v\n", d)
	}
}

func TestMongodb(t *testing.T) {
	var expectErr, actualErr error
	var expectBool, actualBool bool

	// 数据一开始不存在
	expectBool = false
	actualBool = mongodb.Exists(d)
	assert.Equal(t, expectBool, actualBool)

	// 不存在时查找
	expectErr = fmt.Errorf("mongo: no documents in result")
	actualErr = mongodb.FindOne(d)
	assert.Equal(t, expectErr, actualErr)

	// 插入数据, 空值也是值
	expectErr = nil
	actualErr = mongodb.InsertOne(d)
	assert.Equal(t, expectErr, actualErr)

	// 插入后存在
	expectBool = true
	actualBool = mongodb.Exists(d)
	assert.Equal(t, expectBool, actualBool)

	// 根据主键更新数据
	d.URL = "测试网址"
	expectErr = nil
	actualErr = mongodb.UpdataOne(d)
	assert.Equal(t, expectErr, actualErr)

	// 存在时查找
	expectErr = nil
	actualErr = mongodb.FindOne(dd)
	assert.Equal(t, expectErr, actualErr)
	assert.Equal(t, d, dd)

	// 删除数据
	expectErr = nil
	actualErr = mongodb.DeleteOne(dd)
	assert.Equal(t, expectErr, actualErr)
}

func TestRedis(t *testing.T) {
	var expectErr, actualErr error
	var expectBool, actualBool bool

	// 数据一开始不存在
	expectBool = false
	actualBool = redis.Exists(d)
	assert.Equal(t, expectBool, actualBool)

	// 不存在时查找
	expectErr = fmt.Errorf("redigo: nil returned")
	actualErr = redis.FindOne(d)
	assert.Equal(t, expectErr, actualErr)

	// 插入数据, 空值也是值
	expectErr = nil
	actualErr = redis.InsertOne(d)
	assert.Equal(t, expectErr, actualErr)

	// 插入后存在
	expectBool = true
	actualBool = redis.Exists(d)
	assert.Equal(t, expectBool, actualBool)

	// 根据主键更新数据
	d.URL = "测试网址"
	expectErr = nil
	actualErr = redis.UpdataOne(d)
	assert.Equal(t, expectErr, actualErr)

	// 存在时查找
	expectErr = nil
	actualErr = redis.FindOne(dd)
	assert.Equal(t, expectErr, actualErr)
	assert.Equal(t, d, dd)

	// 删除数据
	expectErr = nil
	actualErr = redis.DeleteOne(dd)
	assert.Equal(t, expectErr, actualErr)
}

func TestMongodbCRUD(t *testing.T) {
	// 在 mongodb 中的插入, 查询, 删除测试
	tests := []*AbsText{
		NewAbsText("网址1", []string{"关键字1", "关键字2"}),
		NewAbsText("网址2", []string{"关键字1", ""}),
		NewAbsText("网址3", []string{"关键字1"}),
		NewAbsText("网址4", []string{""}),
		NewAbsText("网址5", []string{}),
		NewAbsText("网址6", nil),
	}
	var err error
	newAT := &AbsText{}

	for _, at := range tests {
		// 保存数据
		err = mongodb.InsertOne(at)
		assert.Equal(t, nil, err)

		// 读取数据
		newAT.Hash = at.Hash
		err = mongodb.FindOne(newAT)
		assert.Equal(t, nil, err)
		assert.Equal(t, at, newAT)

		// 删除数据
		err = mongodb.DeleteOne(at)
		assert.Equal(t, nil, err)
	}
}

func TestRedisCRUD(t *testing.T) {
	// 在 mongodb 中的插入, 查询, 删除测试
	tests := []*AbsText{
		NewAbsText("网址1", []string{"关键字1", "关键字2"}),
		NewAbsText("网址2", []string{"关键字1", ""}),
		NewAbsText("网址3", []string{"关键字1"}),
		NewAbsText("网址4", []string{""}),
		NewAbsText("网址5", []string{}),
		NewAbsText("网址6", nil),
	}
	var err error

	// redis 反序列化时, 不会修改空值数据, 因此接收时尽量每次都是一个空的对象(除了主键)
	for _, at := range tests {
		newAT := &AbsText{}
		// 保存数据
		err = redis.InsertOne(at)
		assert.Equal(t, nil, err)

		// 读取数据
		newAT.Hash = at.Hash
		err = redis.FindOne(newAT)
		if len(at.KeyWords) == 0 {
			at.KeyWords = nil
		}
		if len(at.Abstract) == 0 {
			at.Abstract = nil
		}
		assert.Equal(t, nil, err)
		assert.Equal(t, at, newAT)

		// 删除数据
		err = redis.DeleteOne(at)
		assert.Equal(t, nil, err)
	}
}
