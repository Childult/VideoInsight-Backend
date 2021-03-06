package server

import (
	"fmt"
	"os"
	"swc/mongodb"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func Test(t *testing.T) {
	path := "/home/download/1615024706/b'MTYxNTAyNDcxNC4xNzQ0ODAyaHR0cHM6Ly93d3cuYmlsaWJpbGkuY29tL3ZpZGVvL0JWMU1LNHkxRDdpTg=='.mp4"
	x := Exists(path)
	fmt.Println(x)
}

func TestDownload(t *testing.T) {
	mongodb.SWCDB = "test"
	job := mongodb.Job{URL: "https://www.bilibili.com/video/BV1MK4y1D7iN"}
	StartTask(job)

	assert.Equal(t, nil, nil)
	// filter := bson.M{"url": job.URL}
	// mongodb.DeleteOneByfilter("job", filter)
	// mongodb.DeleteOneByfilter("source", filter)
}
