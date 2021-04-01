package server

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"swc/data/job"
	"swc/data/resource"
	"swc/util"
	"sync"
	"testing"
	"time"
)

// 回调完成标志
var sw = sync.WaitGroup{}
var isCompleted = false

func simplifiedHandle(job *job.Job, result []string) {
	// 是否执行成功
	if len(result) == 0 {
		fmt.Println("[python回调] 执行失败")
	} else {
		fmt.Println("[python回调] 执行成功, 结果如下")
		pythonReturn := strings.Join(result, "")
		fmt.Println(pythonReturn)
	}
	isCompleted = true
	sw.Done()
}

func TestPythonCaller(t *testing.T) {
	// 测试 python 调用模块
	// 新建一个测试用的 python 文件, 写入一些简单的函数
	file, err := os.Create("TestPythonCaller.py")
	if err != nil {
		t.Errorf("[python 调用模块测试] 文件创建失败: %s", err)
		return
	}
	file.WriteString("import time;\ndef getTime(*list,**map):\n return list, map, time.strftime('%Y-%m-%d %H:%M:%S',time.localtime())")
	file.Close()

	// 构建视频下载对象
	start := time.Now()

	python := PyWorker{
		PackagePath: "",                 // python 路径, 当前文件夹
		FileName:    "TestPythonCaller", // 文件名
		MethodName:  "getTime",          // 调用函数
		Args: []string{ // 实参, 这里可以输入任意个数
			SetArg(1),   // 数字 1
			SetArg("2"), // 字符串 2
		},
	}
	sw.Add(1)
	go python.Call(nil, simplifiedHandle)
	sw.Wait()
	fmt.Println("[python 调用模块测试] 完成, 共使用", time.Since(start))
	// 删除测试用的 python 文件
	err = os.Remove("TestPythonCaller.py")
	if err != nil {
		t.Errorf("[python 调用模块测试] 文件删除失败: %s", err)
	}
}

func TestVideoDownload(t *testing.T) {
	// 测试视频下载模块
	fmt.Println("[视频下载模块测试] 开始")
	start := time.Now()

	// 构建需要的对象
	r := resource.Resource{
		URL:      "https://www.bilibili.com/video/BV1PK411w7h8", // 测试连接
		Location: "/swc/resource/test",                          // 保存地址
	}
	// 构建视频下载对象
	python := PyWorker{
		PackagePath: filepath.Join(util.WorkSpace, "video_getter"), // python 包地址
		FileName:    "api",                                         // 文件名
		MethodName:  "download_video",                              // 调用函数
		Args: []string{ // 实参
			SetArg(r.URL),      // 资源链接
			SetArg(r.Location), // 保存路径
		},
	}
	sw.Add(1)
	go python.Call(nil, simplifiedHandle)
	sw.Wait()
	fmt.Println("[视频下载模块测试] 完成, 共使用", time.Since(start))
}
