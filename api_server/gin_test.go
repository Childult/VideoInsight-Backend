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

// videoTest 用于存储视频摘要的结果
type videoTest struct {
	VAbstract []string `json:"VAbstract"`
	Error     string   `json:"Error"`
}

// videoHandle 视频分析的回调
func videoTestHandle(job *job.Job, result []string) {
	job.Status = 1
}

func Test(t *testing.T) {
	mongodb.SWCDB = "test"
	logger.InitLog()
	// python 测试
	python := server.PyWorker{
		PackagePath: "/swc/resource/test",
		FileName:    "a",
		MethodName:  "x",
		Args:        []string{},
	}
	jobs := job.Job{
		DeviceID: "json.DeviceID",
		URL:      "json.URL",
		KeyWords: []string{},
		JobID:    "json.GetID()",
		Status:   0,
	}
	python.Call(&jobs, videoTestHandle)
	for {
		if jobs.Status != 1 {
			time.Sleep(time.Second * 10)
		} else {
			break
		}
	}
	fmt.Printf("%+v, %+v\n", python, jobs)
}

func TestVideoAbstract(t *testing.T) {
	mongodb.SWCDB = "test"
	logger.InitLog()
	// python 测试
	python := server.PyWorker{
		PackagePath: "/swc/code/video_analysis/",
		FileName:    "api",
		MethodName:  "generate_abstract_from_video",
		Args: []string{
			// server.SetArg("/swc/code/api_server/test/test_media.mp4"),
			server.SetArg("/swc/resource/test/media/3.mp4"),
			server.SetArg("/swc/resource/test/go/"),
		},
	}
	jobs := job.Job{
		DeviceID: "json.DeviceID",
		URL:      "json.URL",
		KeyWords: []string{},
		JobID:    "json.GetID()",
		Status:   0,
	}
	python.Call(&jobs, videoTestHandle)
	for {
		if jobs.Status != 1 {
			time.Sleep(time.Second * 100)
		} else {
			break
		}
	}
	fmt.Printf("%+v, %+v\n", python, jobs)
}

func TestDeleteAll(t *testing.T) {
	// 删除所有数据
	mongodb.SWCDB = "test"
	logger.InitLog()
	filter := bson.M{}

	err := mongodb.DeleteManyByfilter("job", filter)
	assert.Equal(t, nil, err)

	err = mongodb.DeleteManyByfilter("resource", filter)
	assert.Equal(t, nil, err)

	err = mongodb.DeleteManyByfilter("abstext", filter)
	assert.Equal(t, nil, err)

	err = mongodb.DeleteManyByfilter("absvideo", filter)
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

	data, err = mongodb.FindOneByfilter("absvideo", filter)
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

	// 启动服务器
	router := GinRouter()

	// 测试 POST
	testPOST := struct {
		method string
		url    string
		body   string
		header map[string]string
		status int
		result string
	}{
		method: "POST",
		url:    "/job",
		body:   `{"deviceid":"1", "url":"https://www.bilibili.com/video/BV18r4y1A7Uv"}`,
		header: map[string]string{"Content-Type": "application/json;charset=utf-8"},
		status: 200,
		result: `{"job_id":"b89259c2efa49ae0cd3db665"}`,
	}

	fmt.Println("======================= 开始测试 POST =====================================")
	postRecorder := httptest.NewRecorder()
	req, _ := http.NewRequest(testPOST.method, testPOST.url, strings.NewReader(testPOST.body))
	for key, value := range testPOST.header {
		req.Header.Set(key, value)
	}
	router.ServeHTTP(postRecorder, req)

	assert.Equal(t, testPOST.status, postRecorder.Code)
	assert.Equal(t, testPOST.result, postRecorder.Body.String())

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
		url:    "/job/b89259c2efa49ae0cd3db665",
		body:   `{"deviceid":"1", "url":"https://www.bilibili.com/video/BV18r4y1A7Uv"}`,
		header: map[string]string{"Content-Type": "application/json;charset=utf-8"},
		status: 200,
		result: `{"status":32}`,
	}

	for {
		getRecorder := httptest.NewRecorder()
		req, _ := http.NewRequest(testGET.method, testGET.url, strings.NewReader(testGET.body))
		for key, value := range testGET.header {
			req.Header.Set(key, value)
		}
		router.ServeHTTP(getRecorder, req)

		if getRecorder.Body.String() == testGET.result {
			break
		} else {
			fmt.Println("==============================", getRecorder.Code, getRecorder.Body.String(), "==============================")
			time.Sleep(time.Second * 10)
		}
	}

}
