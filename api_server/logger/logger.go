package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"swc/util"
)

var (
	// Info 调试
	Info *log.Logger
	// Warning 警告
	Warning *log.Logger
	// Debug 常规
	Debug *log.Logger
	// Error 错误
	Error *log.Logger
)

// InitLog 初始化
func InitLog() {
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

	Info = log.New(io.MultiWriter(os.Stderr, infoFile, debugFile), "Debug:", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(io.MultiWriter(os.Stderr, warningFile, debugFile), "Warning:", log.Ldate|log.Ltime|log.Lshortfile)
	Debug = log.New(io.MultiWriter(debugFile), "Info:", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(os.Stderr, errFile, debugFile), "Error:", log.Ldate|log.Ltime|log.Lshortfile)
}
