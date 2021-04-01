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
	InitRedis("172.17.0.4:6379", "")
}

func TestSetGet(t *testing.T) {
	// 简单的插入查询测试
	var expectErr, actualErr error
	conn := Get()      // 获取连接
	defer conn.Close() // 释放连接

	// 切换到1号数据库, 防止数据污染
	expectErr = nil
	_, actualErr = conn.Do("select", 1)
	assert.Equal(t, expectErr, actualErr)

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

	// 切换到1号数据库, 防止数据污染
	expectErr = nil
	_, actualErr = conn.Do("select", 1)
	assert.Equal(t, expectErr, actualErr)

	expectErr = nil
	_, actualErr = conn.Do("flushall")
	assert.Equal(t, expectErr, actualErr)
}

func TestSelect(t *testing.T) {
	// 简单的插入查询测试
	var expectErr, actualErr error
	conn := Get()      // 获取连接
	defer conn.Close() // 释放连接

	// 切换到1号数据库
	expectErr = nil
	_, actualErr = conn.Do("select", 1)
	assert.Equal(t, expectErr, actualErr)

	// 简单插入数据
	expectErr = nil
	_, actualErr = conn.Do("set", "1", "hello")
	assert.Equal(t, expectErr, actualErr)

	expectErr = nil
	_, actualErr = conn.Do("set", "2", "world")
	assert.Equal(t, expectErr, actualErr)

	// 从连接池里再拿一个连接
	conn2 := Get()
	defer conn2.Close()

	// 查询, 发现查不到, 说明 select 只作用于一个连接, 与整个连接池无关
	expectErr = fmt.Errorf("redigo: nil returned")
	value, actualErr := redis.String(conn2.Do("get", "1"))
	assert.Equal(t, "", value)
	assert.Equal(t, expectErr, actualErr)

	// 切换到1号数据库, 查找并删除
	expectErr = nil
	_, actualErr = conn2.Do("select", 1)
	assert.Equal(t, expectErr, actualErr)

	expectErr = nil
	value, actualErr = redis.String(conn2.Do("get", "1"))
	assert.Equal(t, "hello", value)
	assert.Equal(t, expectErr, actualErr)

	expectErr = nil
	value, actualErr = redis.String(conn2.Do("get", "2"))
	assert.Equal(t, "world", value)
	assert.Equal(t, expectErr, actualErr)

	expectErr = nil
	_, actualErr = conn2.Do("del", "1")
	assert.Equal(t, expectErr, actualErr)

	expectErr = nil
	_, actualErr = conn2.Do("del", "2")
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

	// 切换到1号数据库, 防止数据污染
	expectErr = nil
	_, actualErr = conn.Do("select", 1)
	assert.Equal(t, expectErr, actualErr)

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

	// 切换到1号数据库, 防止数据污染
	expectErr = nil
	_, actualErr = conn.Do("select", 1)
	assert.Equal(t, expectErr, actualErr)

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

	// 反序列化读取数据
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

// fakeData 测试数据
type fakeData struct {
	Addr string
	Name string
}

func (r *fakeData) Value() string {
	return r.Addr
}

func TestCRUD(t *testing.T) {
	// 使用序列化的方式插入数据库
	var expectErr, actualErr error
	var expectBool, actualBool bool
	var d *fakeData = new(fakeData)
	conn := Get()      // 获取连接
	defer conn.Close() // 释放连接

	// 切换到1号数据库, 防止数据污染
	expectErr = nil
	_, actualErr = conn.Do("select", 1)
	assert.Equal(t, expectErr, actualErr)

	// 数据一开始不存在
	expectBool = false
	actualBool = Exists(d)
	assert.Equal(t, expectBool, actualBool)

	// 不存在时查找
	expectErr = fmt.Errorf("redigo: nil returned")
	actualErr = FindOne(d)
	assert.Equal(t, expectErr, actualErr)

	// 插入数据, 空值也是值
	expectErr = nil
	actualErr = InsertOne(d)
	assert.Equal(t, expectErr, actualErr)

	// 插入后存在
	expectBool = true
	actualBool = Exists(d)
	assert.Equal(t, expectBool, actualBool)

	// 根据主键更新数据
	d.Name = "鸡排卷饼"
	expectErr = nil
	actualErr = UpdataOne(d)
	assert.Equal(t, expectErr, actualErr)

	// 存在时查找
	var dd *fakeData = new(fakeData)
	expectErr = nil
	actualErr = FindOne(dd)
	assert.Equal(t, expectErr, actualErr)
	assert.Equal(t, d, dd)

	// 删除数据
	expectErr = nil
	actualErr = DeleteOne(dd)
	assert.Equal(t, expectErr, actualErr)
}
