package redis

import (
	"bytes"
	"encoding/gob"

	"github.com/gomodule/redigo/redis"
)

// Exists 通过数据的主键, 查看数据是否存在
func Exists(data redisData) (b bool) {
	conn := Get()                                      // 获取连接
	defer conn.Close()                                 // 释放连接
	b, _ = redis.Bool(conn.Do("exists", data.Value())) // 判断是否存在
	return
}

// InsertOne 通过主键插入数据
func InsertOne(data redisData) (err error) {
	conn := Get()      // 获取连接
	defer conn.Close() // 释放连接

	var buf bytes.Buffer
	encode := gob.NewEncoder(&buf)
	err = encode.Encode(data)
	if err != nil {
		return
	}

	_, err = conn.Do("set", data.Value(), buf.Bytes())
	return
}

// DeleteOne 删除数据
func DeleteOne(data redisData) (err error) {
	conn := Get()      // 获取连接
	defer conn.Close() // 释放连接
	_, err = conn.Do("del", data.Value())
	return
}

// UpdataOne 取代主键相同的数据
func UpdataOne(data redisData) (err error) {
	return InsertOne(data)
}

// FindOne 返回找到的第一个数据, 以主键查找
func FindOne(data redisData) (err error) {
	conn := Get()      // 获取连接
	defer conn.Close() // 释放连接

	readBytes, err := redis.Bytes(conn.Do("get", data.Value()))
	if err != nil {
		return
	}

	reader := bytes.NewReader(readBytes)
	decode := gob.NewDecoder(reader)

	// 反序列化
	err = decode.Decode(data)
	return
}
