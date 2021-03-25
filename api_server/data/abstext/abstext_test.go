package abstext

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
	tests := []*AbsText{
		NewAbsText("网址1", []string{"关键字1", "关键字2"}),
		NewAbsText("网址2", []string{"关键字1", ""}),
		NewAbsText("网址3", []string{"关键字1"}),
		NewAbsText("网址4", []string{""}),
		NewAbsText("网址5", []string{}),
		NewAbsText("网址6", nil),
	}
	var err error
	newAT := &AbsText{}

	for _, at := range tests {
		// 保存数据
		err = at.Dump()
		assert.Equal(t, nil, err)

		// 读取数据
		newAT.Hash = at.Hash
		err = newAT.Load()
		assert.Equal(t, nil, err)
		assert.Equal(t, at, newAT)

		// 删除数据
		err = at.Delete()
		assert.Equal(t, nil, err)
	}
}
