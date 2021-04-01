package absvideo

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
	d  *AbsVideo = new(AbsVideo)
	dd *AbsVideo = new(AbsVideo)
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
	d.Abstract = []string{"hello", "world"}
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
	d.Abstract = []string{"hello", "world"}
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
	tests := []*AbsVideo{
		{"网址1", []string{"摘要1", "摘要2"}},
		{"网址2", []string{"摘要1", ""}},
		{"网址3", []string{"", ""}},
		{"网址4", []string{""}},
		{"网址5", []string{}},
		{"网址6", nil},
	}
	var err error
	newAV := &AbsVideo{}

	for _, av := range tests {
		// 保存数据
		err = mongodb.InsertOne(av)
		assert.Equal(t, nil, err)

		// 读取数据
		newAV.URL = av.URL
		err = mongodb.FindOne(newAV)
		assert.Equal(t, nil, err)
		assert.Equal(t, av, newAV)

		// 删除数据
		err = mongodb.DeleteOne(newAV)
		assert.Equal(t, nil, err)
	}
}

func TestRedisCRUD(t *testing.T) {
	// 在 mongodb 中的插入, 查询, 删除测试
	tests := []*AbsVideo{
		{"网址1", []string{"摘要1", "摘要2"}},
		{"网址2", []string{"摘要1", ""}},
		{"网址3", []string{"", ""}},
		{"网址4", []string{""}},
		{"网址5", []string{}},
		{"网址6", nil},
	}
	var err error

	// redis 反序列化时, 不会修改空值数据, 因此接收时尽量每次都是一个空的对象(除了主键)
	for _, av := range tests {
		newAV := &AbsVideo{}
		// 保存数据
		err = redis.InsertOne(av)
		assert.Equal(t, nil, err)

		// 读取数据
		newAV.URL = av.URL
		err = redis.FindOne(newAV)
		assert.Equal(t, nil, err)
		if len(av.Abstract) == 0 {
			av.Abstract = nil
		}
		assert.Equal(t, av, newAV)

		// 删除数据
		err = redis.DeleteOne(newAV)
		assert.Equal(t, nil, err)
	}
}
