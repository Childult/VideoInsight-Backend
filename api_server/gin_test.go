package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"swc/data/resource"
	"swc/dbs/mongodb"
	"swc/dbs/redis"
	"swc/util"
	"testing"
)

func init() {
	redis.InitRedis(util.RedisAddr, util.RedisPW)                     // 连接 redis
	mongodb.InitMongodb(util.MongoAddr, util.MongoUser, util.MongoPW) // 连接 mongodb
}

func TestGetData(t *testing.T) {
	// 启动服务器
	router := GinRouter()

	// 测试 GET
	testGET := struct {
		method string
		url    string
		body   string
		header map[string]string
		status int
		result string
	}{
		method: "POST",
		url:    "/job",
		body:   `{"deviceid":"12345dsfsfsdfd", "url":"https://www.bilibili.com/video/BV1Li4y1P7ic"}`,
		header: map[string]string{"Content-Type": "application/json;charset=utf-8"},
		status: 200,
		result: `{"status":32}`,
	}

	getRecorder := httptest.NewRecorder()
	req, _ := http.NewRequest(testGET.method, testGET.url, strings.NewReader(testGET.body))
	for key, value := range testGET.header {
		req.Header.Set(key, value)
	}
	router.ServeHTTP(getRecorder, req)
	fmt.Println("标记", getRecorder)
}

func TestGetJobID(t *testing.T) {
	// 启动服务器
	router := GinRouter()

	// 测试 GET
	testGET := struct {
		method string
		url    string
		body   string
		header map[string]string
		status int
		result string
	}{
		method: "GET",
		url:    "/job/795b95bbf2f756d233b1a575",
		body:   `{"deviceid":"1", "url":"https://www.bilibili.com/video/BV18r4y1A7Uv"}`,
		header: map[string]string{"Content-Type": "application/json;charset=utf-8"},
		status: 200,
		result: `{"status":32}`,
	}

	getRecorder := httptest.NewRecorder()
	req, _ := http.NewRequest(testGET.method, testGET.url, strings.NewReader(testGET.body))
	for key, value := range testGET.header {
		req.Header.Set(key, value)
	}
	router.ServeHTTP(getRecorder, req)
	fmt.Println("标记", getRecorder)

}

func TestFind(t *testing.T) {
	r := &resource.Resource{URL: "https://www.bilibili.com/video/BV1Li4y1P7ic"} // 构建资源
	redis.FindOne(r)
	fmt.Println(r)
	mongodb.FindOne(r)
	fmt.Println(r)
	// redis.UpdataOne(r)
}
