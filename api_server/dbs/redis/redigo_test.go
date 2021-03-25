package redis

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
)

func init() {
	InitRedis("172.17.0.5:6379", "")
}

func TestSetGet(t *testing.T) {
	// 简单的插入查询测试
	var expectErr, actualErr error
	conn := Get()      // 获取连接
	defer conn.Close() // 释放连接

	// 简单插入
	expectErr = nil
	_, actualErr = conn.Do("set", "key", "value", "EX", "1") // `EX 1`表示数据1秒后超时(自动清除数据), `set`不区分大小写
	assert.Equal(t, expectErr, actualErr)

	expectErr = nil
	_, actualErr = redis.String(conn.Do("get", "key"))
	assert.Equal(t, expectErr, actualErr)

	time.Sleep(2 * time.Second)
	expectErr = fmt.Errorf("redigo: nil returned")
	_, actualErr = redis.String(conn.Do("get", "key"))
	assert.Equal(t, expectErr, actualErr)
}

func TestClear(t *testing.T) {
	// 清空所有数据库
	var expectErr, actualErr error
	conn := Get()      // 获取连接
	defer conn.Close() // 释放连接

	expectErr = nil
	_, actualErr = conn.Do("flushall")
	assert.Equal(t, expectErr, actualErr)
}

type testStruct1 struct {
	List []string
	Time time.Time
	Id   string
}

type testStruct2 struct {
	Time string `redis:"t"`
	Id   int64  `redis:"id"`
}

func TestStruct(t *testing.T) {
	// 使用自带的方式将结构体插入 redis数据库
	var expectErr, actualErr error
	conn := Get()      // 获取连接
	defer conn.Close() // 释放连接

	ts1 := testStruct1{
		[]string{"hello"},
		time.Now(),
		"test_id",
	}

	// 自带的 redis.Args, AddFlat 只能针对一些简单的类型(Integer, float, boolean, string and []byte), 否则会无法取出
	expectErr = nil
	_, actualErr = conn.Do("hmset", redis.Args{"id1"}.AddFlat(&ts1)...)
	assert.Equal(t, expectErr, actualErr)

	value, actualErr := redis.Values(conn.Do("hgetall", "id1"))
	assert.Equal(t, expectErr, actualErr)

	expectErr = fmt.Errorf("redigo.ScanStruct: cannot assign field List: cannot convert from Redis bulk string to []string")
	actualErr = redis.ScanStruct(value, &testStruct1{})
	assert.Equal(t, expectErr, actualErr)

	ts2 := testStruct2{
		strconv.FormatInt(time.Now().Unix(), 10),
		123,
	}

	// 都是简单类型就可以取出来
	expectErr = nil
	_, actualErr = conn.Do("hmset", redis.Args{"id2"}.AddFlat(&ts2)...)
	assert.Equal(t, expectErr, actualErr)

	value, actualErr = redis.Values(conn.Do("hgetall", "id2"))
	assert.Equal(t, expectErr, actualErr)

	ts3 := testStruct2{}
	actualErr = redis.ScanStruct(value, &ts3)
	assert.Equal(t, expectErr, actualErr)
	assert.Equal(t, ts2, ts3)

	// 清空
	for _, id := range []string{"id1", "id2"} {
		_, actualErr = conn.Do("del", id)
		assert.Equal(t, expectErr, actualErr)
	}
}

type testStruct3 struct {
	List []string  `redis:"list"`
	Time time.Time `redis:"time"`
	Id   int64     `redis:"id"`
}

func TestGobEncode(t *testing.T) {
	// 使用序列化的方式插入数据库
	var expectErr, actualErr error
	conn := Get()      // 获取连接
	defer conn.Close() // 释放连接

	ts1 := testStruct3{
		[]string{"hello", "world"},
		time.Now(),
		123,
	}

	// 使用序列化的方式, 可以插入任意类型, 这里使用gob序列化
	expectErr = nil
	var buf bytes.Buffer
	encode := gob.NewEncoder(&buf)
	actualErr = encode.Encode(ts1)
	assert.Equal(t, expectErr, actualErr)

	expectErr = nil
	_, actualErr = conn.Do("set", "id1", buf.Bytes())
	assert.Equal(t, expectErr, actualErr)

	expectErr = nil
	readBytes, actualErr := redis.Bytes(conn.Do("get", "id1"))
	assert.Equal(t, expectErr, actualErr)

	// 实际中还是出问题了, 基本一致, 但是省略了一些信息, time 序列化前后不一致
	// 可能是因为 time 里面的属性是私有的, 反射时找不到, 因此建议序列化对象的属性都应该是 public
	reader := bytes.NewReader(readBytes)
	decode := gob.NewDecoder(reader)
	ts2 := testStruct3{}

	expectErr = nil
	actualErr = decode.Decode(&ts2)
	assert.Equal(t, expectErr, actualErr)

	assert.Equal(t, ts1.Id, ts2.Id)
	assert.Equal(t, ts1.List, ts2.List)

	expectErr = nil
	_, actualErr = conn.Do("del", "id1")
	assert.Equal(t, expectErr, actualErr)
}
