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
	// Debug 调试
	Debug *log.Logger
	// Error 错误
	Error *log.Logger
)

// InitLog 初始化
func InitLog() {
	errFile, err := os.OpenFile(filepath.Join(util.Location, "errors.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Open LogFile Error：", err)
	}

	log.SetPrefix("[ VideoInsight ]")

	Info = log.New(os.Stdout, "Info:", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, "Warning:", log.Ldate|log.Ltime|log.Lshortfile)
	Debug = log.New(os.Stdout, "Debug:", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(os.Stderr, errFile), "Error:", log.Ldate|log.Ltime|log.Lshortfile)
}
