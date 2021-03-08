package server

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"swc/mongodb/job"
	"swc/util"
)

// PyWorker python worker
type PyWorker struct {
	PackagePath string
	FileName    string
	MethodName  string
	Args        []string
}

// SetArg toString
func SetArg(i interface{}) string {
	var result string
	switch i := i.(type) {
	case string:
		result = fmt.Sprintf("'%s'", string(i))
	case int:
		result = fmt.Sprintf("%s", strconv.Itoa(i))
	}
	return result
}

// PythonHandlerFunc python 回调函数
type PythonHandlerFunc func(job *job.Job, result []string)

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

// Call 采用管道, 可以实时输出
func (py *PyWorker) Call(job *job.Job, handles ...PythonHandlerFunc) {
	cmd := exec.Command("python3", "-c", py.getCmd())
	fmt.Println(cmd.Args)

	// 获取标准输出
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	// 获取标准错误输出
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	go EasyOut(stderr)
	if handles == nil {
		go HandleOut(stdout, job, nil)
	} else {
		for _, handle := range handles {
			go HandleOut(stdout, job, handle)
		}
	}

	err = cmd.Start()
	// 开始调用
	if err != nil {
		panic(err)
	}
	cmd.Wait()
	return
}

// HandleOut 处理标准输出
func HandleOut(r io.Reader, job *job.Job, handles PythonHandlerFunc) {
	var result []string
	Delimiter := "GoTOPythonDelimiter "
	len := len(Delimiter)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		s := scanner.Text()
		fmt.Println(s)
		index := strings.Index(s, Delimiter)
		if index != -1 {
			result = append(result, s[index+len:])
		}
	}
	if handles != nil {
		go handles(job, result)
	}
}

// EasyOut 简单输出
func EasyOut(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
