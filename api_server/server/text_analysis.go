package server

import (
	"fmt"
	"swc/mongodb"
)

func textAnalysis(job *mongodb.Job, source mongodb.Source) {
	// 构建文本分析对象
	videoGetterPath := "/home/backend/SWC-Backend/text_analysis/"
	fileName := "api"
	methodName := "generate_abstract_from_audio"
	args := []PyArgs{
		ArgsTemp(source.Location + source.AudioPath),
	}
	python := PyWorker{
		PackagePath: videoGetterPath,
		FileName:    fileName,
		MethodName:  methodName,
		Args:        args,
	}
	// 文本分析
	result := python.Call()
	if len(result) != 1 {
		job.Status = mongodb.ErrorHappended
		mongodb.Update(job)
		return
	}

	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++")
	fmt.Println(result)
}
