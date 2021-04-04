package server

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"swc/data/task"
	"swc/dbs/mongodb"
	"swc/dbs/redis"
	"swc/util"
	"testing"
	"time"
)

func init() {
	redis.InitRedis("172.17.0.4:6379", "")
	mongodb.InitMongodb("192.168.2.80:27018", "", "")
	util.MongoDB = "test"
	util.RedisDB = 1
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
	result := python.Call()

	// 是否执行成功
	if len(result) == 0 {
		fmt.Println("[python 调用模块测试] 执行失败")
	} else {
		fmt.Println("[python 调用模块测试] 执行成功, 结果如下")
		pythonReturn := strings.Join(result, "")
		fmt.Println(pythonReturn)
	}

	fmt.Println("[python 调用模块测试] 完成, 共使用", time.Since(start))
	// 删除测试用的 python 文件
	err = os.Remove("TestPythonCaller.py")
	if err != nil {
		t.Errorf("[python 调用模块测试] 文件删除失败: %s", err)
	}
}

// 存活证明
func ImAlive(end chan bool, start time.Time) {
	for {
		select {
		case <-time.NewTicker(time.Second * time.Duration(rand.Int63n(20)+1)).C:
			fmt.Println("已经过去了", time.Since(start))
		case x := <-end:
			if x == true {
				return
			}
		}
	}
}

func TestVideoDownload(t *testing.T) {
	URL := "https://www.bilibili.com/video/BV1wp4y1C7Cz" // 测试连接
	// Location := "/swc/resource/test"                     // 保存地址

	// &{https://www.bilibili.com/video/BV1wp4y1C7Cz 4 /swc/resource/1617284325/ MTYxNzI4NDM3NC40OTY5NzU0aHR0cHM6Ly93d3cuYmlsaWJpbGkuY29tL3ZpZGVvL0JWMXdwNHkxQzdDeg==.mp4 MTYxNzI4NDM3NC40OTY5NzU0aHR0cHM6Ly93d3cuYmlsaWJpbGkuY29tL3ZpZGVvL0JWMXdwNHkxQzdDeg==.mp3 }
	// &{30092e642b82bc6e0e60aa82 https://www.bilibili.com/video/BV1wp4y1C7Cz [] 5 }

	// 测试视频下载模块
	fmt.Println("[视频下载模块测试] 开始")
	job := task.NewTask(URL, nil)
	start := time.Now()
	completed := make(chan bool)
	go ImAlive(completed, start)

	r := creatResource(job) // 创建资源
	mediaDownload(job, r)   // 下载视频
	extractAudio(job, r)    // 提取音频
	completed <- true
	fmt.Println(r)
	fmt.Println(job)
	redis.DeleteOne(r)
	redis.DeleteOne(job)
}
