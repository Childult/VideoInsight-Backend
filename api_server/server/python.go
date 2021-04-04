package server

import (
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"swc/logger"
	"swc/util"
	"sync"
)

// PyWorker 用于调用 python
type PyWorker struct {
	PackagePath string
	FileName    string
	MethodName  string
	Args        []string
}

// SetArg 设置 python 调用的参数, 数字<1>转为<"1">, 字符串<"1">转为<'"1"'>
func SetArg(i interface{}) string {
	var result string
	switch i := i.(type) {
	case string:
		result = fmt.Sprintf("'%s'", string(i))
	case int:
		result = strconv.Itoa(i)
	}
	return result
}

// getCmd 执行命令设置, 调用采用`python -c`, 直接在命令行里写一个简单调用的代码
func (py *PyWorker) getCmd() (r string) {
	r = fmt.Sprintf("\n"+
		"import sys\n"+
		"sys.path.append('%s')\n"+
		"sys.path.append('%s')\n"+
		"import %s as worker\n"+
		"result = worker.%s(%s)\n"+
		"print('GoTOPythonDelimiter',result,end='')\n",
		util.WorkSpace, py.PackagePath, py.FileName, py.MethodName, strings.Join(py.Args, ","))
	return
}

// Call 采用管道和`-c`参数, 可以实时输出
func (py *PyWorker) Call() (result []string) {
	cmd := exec.Command("python3", "-u", "-c", py.getCmd())
	logger.Debug.Println("[python]", cmd.Args)
	wg := sync.WaitGroup{}

	// 获取标准输出
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Error.Println("[python]", err)
		return
	}

	// 获取标准错误输出
	stderr, err := cmd.StderrPipe()
	if err != nil {
		logger.Error.Println("[python]", err)
		return
	}

	// 标准错误直接输出
	wg.Add(1)
	go EasyOut(&wg, stderr)

	// 标准输出需要获取结果
	ch := make(chan []string)
	wg.Add(1)
	go HandleOut(&wg, stdout, ch)

	// 开始调用
	err = cmd.Start()
	if err != nil {
		logger.Error.Println("[python]", err)
		return
	}

	result = <-ch
	wg.Wait()
	logger.Debug.Println("[python] 结束", cmd.Args)
	return
}

var buffSize = 1024 * 10

// HandleOut 处理标准输出
func HandleOut(wg *sync.WaitGroup, r io.Reader, ch chan []string) {
	logger.Debug.Println("[HandleOut] 开始.")
	Delimiter := "GoTOPythonDelimiter "
	len := len(Delimiter)
	isFinded := false
	var result []string

	var sb strings.Builder
	buf := make([]byte, buffSize)

	for {
		sb.Reset()
		n, err := r.Read(buf)
		if err != nil {
			if err != io.EOF {
				logger.Debug.Println("[HandleOut] 读取缓冲区异常, err:", err)
			} else {
				logger.Debug.Println("[HandleOut] 读取缓冲区结束.")
			}
			break
		}
		sb.Write(buf[:n])

		s := sb.String()
		logger.Debug.Println("[python]", s)
		if isFinded {
			result = append(result, s)
		} else {
			index := strings.Index(s, Delimiter)
			if index != -1 {
				result = append(result, s[index+len:])
				isFinded = true
			}
		}
	}
	ch <- result
	wg.Done()
	logger.Debug.Println("[HandleOut] 结束.")
}

// EasyOut 简单输出
func EasyOut(wg *sync.WaitGroup, r io.Reader) {
	logger.Debug.Println("[EasyOut] 开始.")
	var sb strings.Builder
	buf := make([]byte, buffSize)

	for {
		sb.Reset()
		n, err := r.Read(buf)
		if err != nil {
			if err != io.EOF {
				logger.Debug.Println("[EasyOut] 读取缓冲区异常, err:", err)
			} else {
				logger.Debug.Println("[EasyOut] 读取缓冲区结束.")
			}
			break
		}
		sb.Write(buf[:n])

		s := sb.String()
		logger.Debug.Println("[python]", s)
	}
	wg.Done()
	logger.Debug.Println("[EasyOut] 结束.")
}
