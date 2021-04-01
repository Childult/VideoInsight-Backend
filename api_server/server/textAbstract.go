package server

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"swc/data/abstext"
	"swc/data/job"
	"swc/data/resource"
	"swc/logger"
	"swc/util"
)

func extractTextAbstract(job *job.Job) {
	logger.Info.Printf("[提取文本摘要] 开始: %+v.\n", job)
	// 判断资源是否已经存在
	at := abstext.NewAbsText(job.URL, job.KeyWords)
	at.ExistInMongodb()
	if at.ExistInMongodb() { // 文本摘要已经提取完成
		if job.Status&util.JobVideoAbstractExtractionDone != 0 { // 如果视频摘要已经提取完成
			job.Status = util.JobCompleted
			job.Save()
			go JobSchedule(job)
		} else {
			job.Status = job.Status | util.JobTextAbstractExtractionDone
			job.Save()
		}
	} else {
		// 不存在就进行文本分析, 否则忽略
		go textAnalysis(job)
	}
}

// textAnalysis 文本分析
// 不判断该资源是否做过文本分析, 直接插入数据库
func textAnalysis(job *job.Job) {
	logger.Info.Printf("[文本分析] 开始: %+v.\n", job)
	// 获取资源信息
	r := resource.Resource{URL: job.URL}
	err := r.Retrieve()
	if err != nil {
		logger.Error.Printf("[文本分析] 获取资源出错: %+v.\n", err)
		job.Status = util.JobErrFailedToFindResource
		job.Save()
		go JobSchedule(job)
		return
	}

	// 构建文本分析对象
	python := PyWorker{
		PackagePath: filepath.Join(util.WorkSpace, "text_analysis"), // python 文件所在的包
		FileName:    "api",                                          // 文件名
		MethodName:  "generate_abstract_from_audio",                 // 调用函数名
		Args: []string{ // 调用实参
			SetArg(filepath.Join(r.Location, r.AudioPath)),
		},
	}

	// 文本分析
	go python.Call(job, textHandle)
}

// textAbstract 用于存储文本分析的结果
type textAbstract struct {
	AText     string   `json:"AText"`
	TAbstract []string `json:"TAbstract"`
	Error     string   `json:"Error"`
}

// textHandle 文本分析的回调
func textHandle(job *job.Job, result []string) {
	logger.Info.Printf("[文本分析回调] 开始: %+v.\n", job)
	// 获取资源信息
	r := resource.Resource{URL: job.URL}
	err := r.Retrieve()
	if err != nil {
		logger.Error.Printf("[文本分析回调] 获取资源出错: %+v.\n", err)
		job.Status = util.JobErrFailedToFindResource
		job.Save()
		go JobSchedule(job)
		return
	}

	if len(result) == 0 { // 未找到结果
		logger.Error.Println("[文本分析回调] 文本分析失败.")
		job.Status = util.JobErrTextAnalysisFailed
		job.Save()
		go JobSchedule(job)
		return
	}
	pythonReturn := strings.Join(result, "")

	// 找到结果, 从返回结果中提取数据
	var text textAbstract
	err = json.Unmarshal([]byte(pythonReturn), &text)
	if err != nil {
		logger.Error.Printf("[文本分析回调] 从文本分析结果中获取JSON失败: %+v.\n", err)
		job.Status = util.JobErrTextAnalysisReadJSONFailed
		job.Save()
		go JobSchedule(job)
		return
	}

	if text.Error != "" {
		logger.Error.Printf("[文本分析回调] 文本分析失败: %+v.\n", text.Error)
		job.Status = util.JobErrTextAnalysisFailed
		job.Save()
		go JobSchedule(job)
		return
	}

	// 文本分析成功, 初始化需要存储的数据
	abstext := abstext.NewAbsText(job.URL, job.KeyWords)
	abstext.Text = text.AText
	abstext.Abstract = text.TAbstract
	abstext.Dump()

	r.AbsText = abstext.Hash
	r.Save()
	job.AbsText = abstext.Hash
	job.Save()

	if job.Status&util.JobVideoAbstractExtractionDone != 0 { // 视频分析已经完成
		job.Status = util.JobCompleted
		job.Save()
	} else {
		job.Status = util.JobTextAbstractExtractionDone
		job.Save()
	}
	go JobSchedule(job)
}
