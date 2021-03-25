package job

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
	tests := []*Job{
		NewJob("设备ID", "网址", []string{"关键词1", "关键词2"}),
		NewJob("设备ID", "网址", []string{"关键词1"}),
		NewJob("设备ID", "网址", []string{""}),
		NewJob("设备ID", "网址", []string{}),
		NewJob("设备ID", "", []string{}),
		NewJob("", "", []string{}),
	}
	var err error
	newJob := &Job{}

	for _, job := range tests {
		// 保存数据
		err = job.Save()
		assert.Equal(t, nil, err)

		// 读取数据
		newJob.JobID = job.JobID
		err = newJob.Retrieve()
		assert.Equal(t, nil, err)
		assert.Equal(t, job, newJob)

		// 删除数据
		err = job.Remove()
		assert.Equal(t, nil, err)
	}
}

func TestMongodb(t *testing.T) {
	// 在 mongodb 中的插入, 查询, 删除测试
	tests := []*Job{
		NewJob("设备ID", "网址", []string{"关键词1", "关键词2"}),
		NewJob("设备ID", "网址", []string{"关键词1"}),
		NewJob("设备ID", "网址", []string{""}),
		NewJob("设备ID", "网址", []string{}),
		NewJob("设备ID", "", []string{}),
		NewJob("", "", []string{}),
	}
	var err error
	newJob := &Job{}

	for _, job := range tests {
		// 保存数据
		err = job.Dump()
		assert.Equal(t, nil, err)

		// 读取数据
		newJob.JobID = job.JobID
		err = newJob.Load()
		assert.Equal(t, nil, err)
		assert.Equal(t, job, newJob)

		// 删除数据
		err = job.Delete()
		assert.Equal(t, nil, err)
	}
}
