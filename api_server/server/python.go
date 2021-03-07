package server

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"swc/util"
)

// PyArgs is a interface for arguments of python program
type PyArgs interface {
	String() string
}

// PyWorker python worker
type PyWorker struct {
	PackagePath string
	FileName    string
	MethodName  string
	Args        []PyArgs
}

func (py *PyWorker) getCmd() (r string) {
	len := len(py.Args)
	args := make([]string, len)
	for i, arg := range py.Args {
		args[i] = arg.String()
	}

	r = fmt.Sprintf("\n"+
		"import sys\n"+
		"sys.path.append('%s')\n"+
		"sys.path.append('%s')\n"+
		"import %s as worker\n"+
		"result = worker.%s(%s)\n"+
		"print('GoTOPythonDelimiter',result,end='')\n",
		util.WorkSpace, py.PackagePath, py.FileName, py.MethodName, strings.Join(args, ","))
	return
}

// CallPython will execute a python program
func CallPython(packagePath string, fileName string, methodName string, args []PyArgs) (result []string) {
	var py PyWorker = PyWorker{PackagePath: packagePath, FileName: fileName, MethodName: methodName, Args: args}
	return py.Call()
}

// Call 采用管道, 可以实时输出
func (py *PyWorker) Call() (result []string) {
	cmd := exec.Command("python3", "-c", py.getCmd())
	fmt.Println(cmd.Args)

	// 获取标准输出

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	result = Stdout(stdout)
	Stdout(stderr)
	cmd.Wait()
	return
}

// ReCall 采用管道, 可以实时输出
func (py *PyWorker) ReCall(handles ...PythonHandlerFunc) {
	cmd := exec.Command("python3", "-c", py.getCmd())
	fmt.Println(cmd.Args)

	// 获取标准输出

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	for _, handle := range handles {
		go ReStdout(stdout, handle)
	}
	for _, handle := range handles {
		go ReStdout(stderr, handle)
	}
	cmd.Wait()
	return
}

// ArgsTemp is a demo type
type ArgsTemp string

// String toString
func (i ArgsTemp) String() string {
	return fmt.Sprintf("'%s'", string(i))
}

// Stdout 处理标准输出
// func Stdout(r io.Reader) (result []string) {
// 	scanner := bufio.NewScanner(r)
// 	for scanner.Scan() {
// 		result = append(result, scanner.Text())
// 	}
// 	return
// }

// Stdout 处理标准输出
func Stdout(r io.Reader) (result []string) {
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
	return
}

// ReStdout 处理标准输出
func ReStdout(r io.Reader, handles PythonHandlerFunc) {
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
}
