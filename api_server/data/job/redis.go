package job

import (
	"bytes"
	"encoding/gob"
	"fmt"
	swcredis "swc/dbs/redis"

	"github.com/gomodule/redigo/redis"
)

// Save 保存到 redis 中
func (j *Job) Save() (err error) {
	conn := swcredis.Get() // 获取连接
	defer conn.Close()     // 释放连接

	var buf bytes.Buffer
	encode := gob.NewEncoder(&buf)
	err = encode.Encode(*j)
	if err != nil {
		err = fmt.Errorf("<%v>序列化失败; 原始错误<%s>", *j, err)
		return
	}

	_, err = conn.Do("set", j.GetKeyValue(), buf.Bytes())
	if err != nil {
		err = fmt.Errorf("<%v>保存到 redis 失败; 原始错误<%s>", *j, err)
		return
	}
	return
}

// Get 从 redis 中读出
func (j *Job) Retrieve() (err error) {
	conn := swcredis.Get() // 获取连接
	defer conn.Close()     // 释放连接

	readBytes, err := redis.Bytes(conn.Do("get", j.GetKeyValue()))
	if err != nil {
		err = fmt.Errorf("从redis中读取<%v>失败; 原始错误<%s>", j.GetKeyValue(), err)
		return err
	}

	reader := bytes.NewReader(readBytes)
	decode := gob.NewDecoder(reader)

	// 反序列化
	job := Job{}
	err = decode.Decode(&job)
	*j = job
	if err != nil {
		err = fmt.Errorf("反序列化<%v>失败; 原始错误<%s>", j.GetKeyValue(), err)
	}

	return err
}

// Remove 从 redis 中移除
func (j *Job) Remove() (err error) {
	conn := swcredis.Get() // 获取连接
	defer conn.Close()     // 释放连接

	_, err = conn.Do("del", j.GetKeyValue())
	if err != nil {
		err = fmt.Errorf("<%v>移除; 原始错误<%s>", *j, err)
		return err
	}
	return
}
