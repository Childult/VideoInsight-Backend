package absvideo

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

func TestMongodb(t *testing.T) {
	// 在 mongodb 中的插入, 查询, 删除测试
	tests := []*AbsVideo{
		{"网址1", []string{"摘要1", "摘要2"}},
		{"网址2", []string{"摘要1", ""}},
		{"网址3", []string{"", ""}},
		{"网址4", []string{""}},
		{"网址5", []string{}},
		{"网址6", nil},
	}
	var err error
	newAV := &AbsVideo{}

	for _, av := range tests {
		// 保存数据
		err = av.Dump()
		assert.Equal(t, nil, err)

		// 读取数据
		newAV.URL = av.URL
		err = newAV.Load()
		assert.Equal(t, nil, err)
		assert.Equal(t, av, newAV)

		// 删除数据
		err = av.Delete()
		assert.Equal(t, nil, err)
	}
}
