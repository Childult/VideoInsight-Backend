package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"swc/mongodb"
	"swc/server"
	"swc/util"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

// 下载测试
func TestMedia(t *testing.T) {

}

func Test(t *testing.T) {
	filePath := "/home/download/1615035748//home/download/1615035748/MTYxNTAzNTc1NC41MDE4OThodHRwczovL3d3dy5iaWxpYmlsaS5jb20vdmlkZW8vQlYxQUs0eTFKN1pl.mp3"
	result := Exists(filePath)
	fmt.Println(result)
}

func TestDownload(t *testing.T) {
	mongodb.SWCDB = "test"

	mediaURL := "https://www.bilibili.com/video/BV1py4y1a7tP"
	savePath := "/home/download/1234567890/"

	// 构建视频下载对象
	videoGetterPath := filepath.Join(util.WorkSpace, "video_getter")
	fileName := "main"
	methodName := "download_video"
	args := []server.PyArgs{
		server.ArgsTemp(mediaURL),
		server.ArgsTemp(savePath),
	}
	python := server.PyWorker{
		PackagePath: videoGetterPath,
		FileName:    fileName,
		MethodName:  methodName,
		Args:        args,
	}
	python.Call()
}

func TestTextAnalysis(t *testing.T) {
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
			body:   `{"deviceid":"1", "url":"https://www.bilibili.com/video/BV18r4y1A7Uv"}`,
			header: map[string]string{"Content-Type": "application/json;charset=utf-8"},
			status: 200,
			result: `{"jobid":"[241 224 136 245 8 214 233 235 202 207 19 254]"}`,
		},
		{
			method: "DELETE",
			url:    "/job",
			body:   `{"deviceid":"1", "url":"https://www.bilibili.com/video/BV18r4y1A7Uv"}`,
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
	time.Sleep(time.Second * 90)
}
