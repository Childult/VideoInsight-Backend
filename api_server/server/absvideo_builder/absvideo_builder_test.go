package absvideo_builder

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
	util.GRPCAddress = "192.168.2.80:50051" // gRPC 调用地址
}

func TestRequestVideoAnalysis(t *testing.T) {
	type args struct {
		url      string
		path     string
		location string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "视频分析测试", args: args{url: "https://www.bilibili.com/video/BV1cK4y1K7Qy", path: "/swc/resource/1617532055/MTYxNzUzMjA5My44MDUzNDU4aHR0cHM6Ly93d3cuYmlsaWJpbGkuY29tL3ZpZGVvL0JWMWNLNHkxSzdReQ==.mp4", location: "/swc/resource/1617532055/"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RequestVideoAnalysis(tt.args.url, tt.args.path, tt.args.location); (err != nil) != tt.wantErr {
				t.Errorf("RequestVideoAnalysis() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
