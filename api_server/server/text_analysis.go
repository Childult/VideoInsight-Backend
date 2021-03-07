package server

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"swc/mongodb"
	"swc/mongodb/job"
	"swc/mongodb/source"
	"swc/util"
)

type textAbstract struct {
	AText     string `json:"AText"`
	TAbstract string `json:"TAbstract"`
	Error     string `json:"Error"`
}

func textAnalysis(job *job.Job, source source.Source) {
	// 构建文本分析对象
	videoGetterPath := "/home/backend/SWC-Backend/text_analysis/"
	fileName := "api"
	methodName := "generate_abstract_from_audio"
	args := []PyArgs{
		ArgsTemp(filepath.Join(source.Location, source.AudioPath)),
	}

	fmt.Println("======================")
	fmt.Println(source.Location)
	fmt.Println(source.AudioPath)
	python := PyWorker{
		PackagePath: videoGetterPath,
		FileName:    fileName,
		MethodName:  methodName,
		Args:        args,
	}
	// 文本分析
	result := python.Call()
	if len(result) != 1 {
		job.Status = util.ErrorHappended
		mongodb.Update(job)
		return
	}

	var text textAbstract
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++")
	err := json.Unmarshal([]byte(result[0]), &text)
	if err != nil {

	}
	fmt.Printf("%+v\n", text)
}
