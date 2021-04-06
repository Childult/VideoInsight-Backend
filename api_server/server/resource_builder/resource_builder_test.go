package resource_builder

import (
	"swc/dbs/mongodb"
	"swc/dbs/redis"
	"swc/util"
	"testing"
)

func init() {
	redis.InitRedis("172.17.0.3:6379", "")
	mongodb.InitMongodb("192.168.2.80:27018", "", "")
	util.MongoDB = "test"
}

func TestRequestResource(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "下载测试", args: args{"https://www.bilibili.com/video/BV1cK4y1K7Qy"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RequestResource(tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("RequestResource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
