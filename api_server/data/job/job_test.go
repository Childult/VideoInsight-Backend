package job

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
	d *Job = new(Job)
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

	// 删除数据
	expectErr = nil
	actualErr = mongodb.DeleteOne(d)
	assert.Equal(t, expectErr, actualErr)
}

func TestRedisCRUD(t *testing.T) {
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

	// 删除数据
	expectErr = nil
	actualErr = redis.DeleteOne(d)
	assert.Equal(t, expectErr, actualErr)
}
func TestMongodbCRUD(t *testing.T) {
	// 在 mongodb 中的插入, 查询, 删除测试
	tests := []*Job{
		NewJob("id1", "网址", []string{"关键词1", "关键词2"}),
		NewJob("id1", "网址", []string{"关键词1"}),
		NewJob("id1", "网址", []string{""}),
		NewJob("id1", "网址", []string{}),
		NewJob("id1", "", []string{}),
		NewJob("id1", "", []string{}),
		NewJob("id1", "", nil),
	}
	var err error
	newJob := &Job{}

	for _, job := range tests {
		// 保存数据
		err = mongodb.InsertOne(job)
		assert.Equal(t, nil, err)

		// 读取数据
		newJob.JobID = job.JobID
		err = mongodb.FindOne(newJob)
		assert.Equal(t, nil, err)
		assert.Equal(t, job, newJob)

		// 删除数据
		err = mongodb.DeleteOne(job)
		assert.Equal(t, nil, err)
	}
}

func TestRedis(t *testing.T) {
	// 在 redis 中的插入, 查询, 删除测试
	tests := []*Job{
		NewJob("id1", "网址", []string{"关键词1", "关键词2"}),
		NewJob("id1", "网址", []string{"关键词1"}),
		NewJob("id1", "网址", []string{""}),
		NewJob("id1", "网址", []string{}),
		NewJob("id1", "", []string{}),
		NewJob("id1", "", []string{}),
		NewJob("id1", "", nil),
	}
	var err error

	// redis 反序列化时, 不会修改空值数据, 因此接收时尽量每次都是一个空的对象(除了主键)
	for _, job := range tests {
		newJob := &Job{}
		// 保存数据
		err = redis.InsertOne(job)
		assert.Equal(t, nil, err)

		// 读取数据
		newJob.JobID = job.JobID
		err = redis.FindOne(newJob)
		assert.Equal(t, nil, err)
		assert.Equal(t, job, newJob)

		// 删除数据
		err = redis.DeleteOne(newJob)
		assert.Equal(t, nil, err)
	}
}
