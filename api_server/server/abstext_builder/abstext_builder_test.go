package abstext_builder

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

func TestRequestTextAnalysis(t *testing.T) {
	type args struct {
		url      string
		keyWords []string
		path     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "文本分析测试", args: args{url: "https://www.bilibili.com/video/BV1cK4y1K7Qy", keyWords: nil, path: "/swc/resource/1617532055/MTYxNzUzMjA5My44MDUzNDU4aHR0cHM6Ly93d3cuYmlsaWJpbGkuY29tL3ZpZGVvL0JWMWNLNHkxSzdReQ==.mp3"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RequestTextAnalysis(tt.args.url, tt.args.keyWords, tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("RequestTextAnalysis() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
