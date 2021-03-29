package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"swc/logger"
	"swc/mongodb"
	"swc/mongodb/abstext"
	"swc/mongodb/absvideo"
	"swc/mongodb/job"
	"swc/mongodb/resource"
	"swc/server"
	"testing"
	"time"

	pb "swc/server/network"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc"
)

func TestPicture(t *testing.T) {
	mongodb.SWCDB = "test"
	logger.InitLog()

	job, _ := job.GetByKey("b89259c2efa49ae0cd3db665")
	av := absvideo.AbsVideo{URL: job.URL}
	video, _ := mongodb.FindOne(av)
	r, _ := resource.GetByKey("https://www.bilibili.com/video/BV18r4y1A7Uv")

	prefix := r.Location
	pictures := video["abstract"]
	pics := make(map[string]string)

	fmt.Println(prefix)
	for _, x := range pictures.(bson.A) {
		file, _ := os.Open(prefix + x.(string))
		content, _ := ioutil.ReadAll(file)
		pics[x.(string)] = string(content)
	}
	a := pictures.(bson.A)
	fmt.Println(pics[a[1].(string)])

}

func TestGet(t *testing.T) {
	mongodb.SWCDB = "test"
	logger.InitLog()

	// 启动服务器
	router := GinRouter()

	// 测试 Get
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

func TestGetText(t *testing.T) {
	mongodb.SWCDB = "test"
	logger.InitLog()
	at := abstext.NewAbsText("https://www.bilibili.com/video/BV18r4y1A7Uv", "", nil, []string{})

	data, _ := mongodb.FindOne(at)
	for key, value := range data {
		fmt.Println(key, value)
	}
}

func TestGRPC(t *testing.T) {
	// address := "localhost:50051"
	address := "192.168.2.80:50051"
	jobID := "12306"
	filePath := "/swc/code/video_analysis/dataset/1.mp4"
	savaPath := "/swc/resource/"

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewVideoAnalysisClient(conn)

	// Contact the server and print out its response.
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()
	r, err := c.GetStaticVideoAbstract(context.TODO(), &pb.VideoInfo{JobId: jobID, File: filePath, SaveDir: savaPath})
	if err != nil {
		log.Fatalf("could not call: %v", err)
	}
	log.Println("Get: ", r.GetJobID(), r.GetPicName(), r.GetError())
}

// videoHandle 视频分析的回调
func videoTestHandle(job *job.Job, result []string) {
	job.Status = 1
	fmt.Println("回调:", strings.Join(result, ""))
}

func TestMediaDownload(t *testing.T) {
	mongodb.SWCDB = "test"
	logger.InitLog()
	// python 下载
	python := server.PyWorker{
		PackagePath: "/swc/code/video_getter/",
		FileName:    "api",
		MethodName:  "download_video",
		Args: []string{
			// server.SetArg("https://www.bilibili.com/video/BV1st411g7a1"),
			server.SetArg("https://www.bilibili.com/video/BV11y4y177hR"), // 短视频
			server.SetArg("/swc/resource/test/"),
		},
	}
	jobs := job.Job{
		DeviceID: "json.DeviceID",
		URL:      "json.URL",
		KeyWords: []string{},
		JobID:    "json.GetID()",
		Status:   0,
	}
	go python.Call(&jobs, videoTestHandle)
	for {
		if jobs.Status != 1 {
			time.Sleep(time.Second * 10)
			fmt.Println("hello")
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
			// server.SetArg("/swc/resource/test/MTYxNTkwNTU1OS4zNzMyNTM4aHR0cHM6Ly93d3cuYmlsaWJpbGkuY29tL3ZpZGVvL0JWMXN0NDExZzdhMQ==.mp4"),
			server.SetArg("/swc/resource/test/MTYxNTkzODYzMy40ODY0NjA3aHR0cHM6Ly93d3cuYmlsaWJpbGkuY29tL3ZpZGVvL0JWMTF5NHkxNzdoUg==.mp4"), // 短视频
			server.SetArg("/swc/resource/test/"),
		},
	}
	jobs := job.Job{
		DeviceID: "json.DeviceID",
		URL:      "json.URL",
		KeyWords: []string{},
		JobID:    "json.GetID()",
		Status:   0,
	}
	go python.Call(&jobs, videoTestHandle)
	for {
		if jobs.Status != 1 {
			time.Sleep(time.Second * 10)
			fmt.Println("hello")
		} else {
			break
		}
	}
	fmt.Printf("结果: %+v, %+v\n", python, jobs)
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

	err = mongodb.DeleteManyByfilter("abstract_text", filter)
	assert.Equal(t, nil, err)

	err = mongodb.DeleteManyByfilter("abstract_video", filter)
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

	data, err = mongodb.FindOneByfilter("abstract_text", filter)
	expectErr = fmt.Errorf("Not Found <%+v>", filter)
	expectdata = bson.M(nil)
	assert.Equal(t, expectErr, err)
	assert.Equal(t, expectdata, data)

	data, err = mongodb.FindOneByfilter("abstract_video", filter)
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
