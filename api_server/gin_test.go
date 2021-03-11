package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"swc/logger"
	"swc/mongodb"
	"swc/mongodb/job"
	"swc/mongodb/resource"
	"swc/server"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestVideoAbstract(t *testing.T) {
	mongodb.SWCDB = "test"
	logger.InitLog()
	// python 测试
	python := server.PyWorker{
		PackagePath: "/home/backend/SWC-Backend/video_analysis/",
		FileName:    "api",
		MethodName:  "generate_abstract_from_video",
		Args: []string{
			server.SetArg("/home/download/123.mp4"),
			server.SetArg("/home/download/"),
		},
	}
	jobs := job.Job{
		DeviceID: "json.DeviceID",
		URL:      "json.URL",
		KeyWords: []string{},
		JobID:    "json.GetID()",
		Status:   0,
	}
	python.Call(&jobs)
	fmt.Printf("%+v, %+v\n", python, jobs)
}

func TestDeleteAll(t *testing.T) {
	// 删除所有数据
	mongodb.SWCDB = "test"
	filter := bson.M{}

	err := mongodb.DeleteManyByfilter("job", filter)
	assert.Equal(t, nil, err)

	err = mongodb.DeleteManyByfilter("resource", filter)
	assert.Equal(t, nil, err)

	err = mongodb.DeleteManyByfilter("abstext", filter)
	assert.Equal(t, nil, err)

	data, err := mongodb.FindOneByfilter("job", filter)
	expectErr := fmt.Errorf("Not Found <%+v>", filter)
	expectdata := bson.M(nil)
	assert.Equal(t, expectErr, err)
	assert.Equal(t, expectdata, data)

	data, err = mongodb.FindOneByfilter("resource", filter)
	expectErr = fmt.Errorf("Not Found <%+v>", filter)
	expectdata = bson.M(nil)
	assert.Equal(t, expectErr, err)
	assert.Equal(t, expectdata, data)

	data, err = mongodb.FindOneByfilter("abstext", filter)
	expectErr = fmt.Errorf("Not Found <%+v>", filter)
	expectdata = bson.M(nil)
	assert.Equal(t, expectErr, err)
	assert.Equal(t, expectdata, data)
}

func TestResourceInsertDelete(t *testing.T) {
	// reource 的插入删除测试
	mongodb.SWCDB = "test"

	url := "www.baidu.com"
	resource1 := &resource.Resource{
		URL:    url,
		Status: 0,
	}

	mongodb.InsertOne(resource1) // 插入
	resource1.SetStatus(1)       // 状态更新
	assert.Equal(t, int32(1), resource1.Status)

	resource2, err := resource.GetByKey("www.baidu.com") // 读取
	assert.Equal(t, int32(1), resource2.Status)
	assert.Equal(t, nil, err)

	err = mongodb.DeleteOne(resource2) // 删除
	assert.Equal(t, nil, err)

	resource3, err := resource.GetByKey("www.baidu.com") // 删除后读取
	filter := bson.M{"url": url}
	expectErr := fmt.Errorf("Not Found <%+v>", filter)
	assert.Equal(t, expectErr, err)
	assert.Equal(t, "", resource3.AbsText)
}

func TestJob(t *testing.T) {
	// job 读取测试
	mongodb.SWCDB = "test"

	// 切片对象是一个指针, 未赋初值时为 nil, 与 [] 不同. 前者是空指针, 后者是空切片, 二者都可以调用方法.
	job1 := job.Job{
		DeviceID: "1",
	}
	fmt.Printf("%+v\n", job1)
	fmt.Println(`job1["key_words"] == nil? true`)
	assert.Equal(t, true, job1.KeyWords == nil)
	mongodb.InsertOne(job1)

	filter := bson.M{}
	data, _ := mongodb.FindOneByfilter("job", filter)
	fmt.Println(`data["key_words"] == nil? true`)
	assert.Equal(t, true, data["key_words"] == nil)
	mongodb.DeleteOneByfilter("job", filter)

	job2 := job.Job{
		DeviceID: "1",
		KeyWords: []string{},
	}
	fmt.Printf("%+v\n", job2)
	fmt.Println(`job2["key_words"] == nil? false`)
	assert.Equal(t, false, job2.KeyWords == nil)
	mongodb.InsertOne(job2)

	data, _ = mongodb.FindOneByfilter("job", filter)
	fmt.Println(`data["key_words"] == nil? false`)
	assert.Equal(t, false, data["key_words"] == nil)
	mongodb.DeleteOneByfilter("job", filter)

}

func TestTextAnalysis(t *testing.T) {
	mongodb.SWCDB = "test"
	logger.InitLog()

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
			body:   `{"deviceid":"1", "url":"https://www.bilibili.com/video/BV18r4y1A7Uv"}`,
			header: map[string]string{"Content-Type": "application/json;charset=utf-8"},
			status: 200,
			result: `{"jobid":"[184 146 89 194 239 164 154 224 205 61 182 101]"}`,
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
	time.Sleep(time.Second * 90)
}
