package resource

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
	d  *Resource = new(Resource)
	dd *Resource = new(Resource)
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
	d.Location = "/摘要"
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
	d.Location = "/摘要"
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
	// 在 redis 中的插入, 查询, 删除测试
	tests := []*Resource{
		{"网址1", 0, "存储路径", "视频", "音频"},
		{"网址2", 0, "存储路径", "视频", "音频"},
		{"网址3", 0, "存储路径", "视频", ""},
		{"网址4", 0, "存储路径", "", ""},
		{"网址5", 0, "", "", ""},
		{"", 0, "", "", ""},
	}
	var err error
	// mongodb 查找时, 会覆盖数据中原有的内容, 即使找到的值是空
	newResource := &Resource{}

	for _, resource := range tests {
		// 保存数据
		err = mongodb.InsertOne(resource)
		assert.Equal(t, nil, err)

		// 读取数据
		newResource.URL = resource.URL
		err = mongodb.FindOne(newResource)
		assert.Equal(t, nil, err)
		assert.Equal(t, resource, newResource)

		// 删除数据
		err = mongodb.DeleteOne(newResource)
		assert.Equal(t, nil, err)
	}
}

func TestRedisCRUD(t *testing.T) {
	// 在 mongodb 中的插入, 查询, 删除测试
	tests := []*Resource{
		{"网址1", 0, "存储路径", "视频", "音频"},
		{"网址2", 0, "存储路径", "视频", "音频"},
		{"网址3", 0, "存储路径", "视频", ""},
		{"网址4", 0, "存储路径", "", ""},
		{"网址5", 0, "", "", ""},
		{"", 0, "", "", ""},
	}
	var err error

	// redis 反序列化时, 不会修改空值数据, 因此接收时尽量每次都是一个空的对象(除了主键)
	for _, resource := range tests {
		newResource := &Resource{}
		// 保存数据
		err = redis.InsertOne(resource)
		assert.Equal(t, nil, err)

		// 读取数据
		newResource.URL = resource.URL
		err = redis.FindOne(newResource)
		assert.Equal(t, nil, err)
		assert.Equal(t, resource, newResource)

		// 删除数据
		err = redis.DeleteOne(newResource)
		assert.Equal(t, nil, err)
	}
}
