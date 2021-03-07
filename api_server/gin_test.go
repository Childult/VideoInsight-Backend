package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"swc/mongodb"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	dir := filepath.Join("/home/", "download", "a", "b")
	fmt.Println(dir)
	dir, _ = filepath.Split("/home/download/a.map3")
	fmt.Println(dir)
	dir, _ = filepath.Split("/home/download/")
	fmt.Println(dir)
	GinRouter()
}

func TestDownload(t *testing.T) {
	mongodb.SWCDB = "test"
	// job := mongodb.Job{URL: "https://www.bilibili.com/video/BV1MK4y1D7iN"}
	// StartTask(job)
	//
	// filter := bson.M{"url": job.URL}
	// data, err := mongodb.FindOneByfilter("job", filter)
	//
	// assert.Equal(t, nil, err)
	// assert.Equal(t, nil, data)
	// // mongodb.DeleteOneByfilter("job", filter)
	// // mongodb.DeleteOneByfilter("source", filter)

	mongodb.SWCDB = "test"
	tests := []struct {
		method string
		url    string
		body   string
		header map[string]string
		status int
		result string
	}{
		{
			method: "POST",
			url:    "/job",
			body:   `{"deviceid":"1", "url":"https://www.bilibili.com/video/BV1AK4y1J7Ze"}`,
			header: map[string]string{"Content-Type": "application/json;charset=utf-8"},
			status: 200,
			result: `{"jobid":"[241 224 136 245 8 214 233 235 202 207 19 254]"}`,
		},
		{
			method: "DELETE",
			url:    "/job",
			body:   `{"deviceid":"1", "url":"https://www.bilibili.com/video/BV1AK4y1J7Ze"}`,
			header: map[string]string{"Content-Type": "application/json;charset=utf-8"},
			status: 200,
			result: "",
		},
	}

	// 启动服务器
	router := GinRouter()

	for index, test := range tests {
		fmt.Println("====================== 开始测试第", index, "例测试数据 ==================================")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(test.method, test.url, strings.NewReader(test.body))
		for key, value := range test.header {
			req.Header.Set(key, value)
		}
		router.ServeHTTP(w, req)

		assert.Equal(t, test.status, w.Code)
		assert.Equal(t, test.result, w.Body.String())
	}

	// 等待输出
	time.Sleep(time.Second * 10)
}
