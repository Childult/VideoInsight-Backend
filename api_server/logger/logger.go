package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"swc/util"
)

var (
	// Info 常规
	Info *log.Logger
	// Warning 警告
	Warning *log.Logger
	// Debug 调试, 所有信息都会出现在 debug 文件中
	Debug *log.Logger
	// Error 错误
	Error *log.Logger
)

// 初始化日志
func init() {
	infoFile, err := os.OpenFile(filepath.Join(util.LogFile, "info.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Open LogFile Error：", err)
	}
	warningFile, err := os.OpenFile(filepath.Join(util.LogFile, "warning.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Open LogFile Error：", err)
	}
	debugFile, err := os.OpenFile(filepath.Join(util.LogFile, "debug.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Open LogFile Error：", err)
	}
	errFile, err := os.OpenFile(filepath.Join(util.LogFile, "errors.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Open LogFile Error：", err)
	}

	// info, warning, error 对应三个日志等级, 所有消息都会倒到 debug 中
	Info = log.New(io.MultiWriter(os.Stderr, infoFile, debugFile), "Info:", log.LstdFlags|log.Lshortfile)
	Warning = log.New(io.MultiWriter(os.Stderr, warningFile, debugFile), "Warning:", log.LstdFlags|log.Lshortfile)
	Debug = log.New(io.MultiWriter(os.Stderr, debugFile), "Debug:", log.LstdFlags|log.Lshortfile)
	Error = log.New(io.MultiWriter(os.Stderr, errFile, debugFile), "Error:", log.LstdFlags|log.Lshortfile)
}
