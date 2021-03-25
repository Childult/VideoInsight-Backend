package resource

import (
	mymongo "swc/dbs/mongodb"
	myredis "swc/dbs/redis"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	myredis.InitRedis("172.17.0.5:6379", "")
	mymongo.InitMongodb("172.17.0.3:27017", "", "")
}

func TestRedis(t *testing.T) {
	// 在 redis 中的插入, 查询, 删除测试
	tests := []*Resource{
		{"网址1", 0, "存储路径", "视频", "音频", "文本摘要"},
		{"网址2", 0, "存储路径", "视频", "音频", ""},
		{"网址3", 0, "存储路径", "视频", "", ""},
		{"网址4", 0, "存储路径", "", "", ""},
		{"网址5", 0, "", "", "", ""},
		{"", 0, "", "", "", ""},
	}
	var err error
	newResource := &Resource{}

	for _, resource := range tests {
		// 保存数据
		err = resource.Save()
		assert.Equal(t, nil, err)

		// 读取数据
		newResource.URL = resource.URL
		err = newResource.Retrieve()
		assert.Equal(t, nil, err)
		assert.Equal(t, resource, newResource)

		// 删除数据
		err = resource.Remove()
		assert.Equal(t, nil, err)
	}
}

func TestMongodb(t *testing.T) {
	// 在 mongodb 中的插入, 查询, 删除测试
	tests := []*Resource{
		{"网址1", 0, "存储路径", "视频", "音频", "文本摘要"},
		{"网址2", 0, "存储路径", "视频", "音频", ""},
		{"网址3", 0, "存储路径", "视频", "", ""},
		{"网址4", 0, "存储路径", "", "", ""},
		{"网址5", 0, "", "", "", ""},
		{"", 0, "", "", "", ""},
	}
	var err error
	newResource := &Resource{}

	for _, resource := range tests {
		// 保存数据
		err = resource.Dump()
		assert.Equal(t, nil, err)

		// 读取数据
		newResource.URL = resource.URL
		err = newResource.Load()
		assert.Equal(t, nil, err)
		assert.Equal(t, resource, newResource)

		// 删除数据
		err = resource.Delete()
		assert.Equal(t, nil, err)
	}
}
