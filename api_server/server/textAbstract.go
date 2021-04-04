package server

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"swc/data/abstext"
	"swc/data/resource"
	"swc/data/task"
	"swc/dbs/mongodb"
	"swc/dbs/redis"
	"swc/logger"
	"swc/util"
)

// textAbstract 用于存储文本分析的结果
type textAbstract struct {
	AText     string   `json:"AText"`
	TAbstract []string `json:"TAbstract"`
	Error     string   `json:"Error"`
}

func extractTextAbstract(task *task.Task, r *resource.Resource) {
	logger.Info.Printf("[提取文本摘要] 开始: %+v.\n", task)
	// 开始进行文本分析

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
	result := python.Call()

	if len(result) == 0 { // 未找到结果
		logger.Error.Println("[提取文本摘要] 文本分析失败.")
		task.Status = util.JobErrTextAnalysisFailed
		redis.UpdataOne(task)
		return
	}

	pythonReturn := strings.Join(result, "")

	// 找到结果, 从返回结果中提取数据
	var text textAbstract
	err := json.Unmarshal([]byte(pythonReturn), &text)
	if err != nil {
		logger.Error.Printf("[提取文本摘要] 从文本分析结果中获取JSON失败: %+v.\n", err)
		task.Status = util.JobErrTextAnalysisReadJSONFailed
		redis.UpdataOne(task)
		return
	}

	if text.Error != "" {
		logger.Error.Printf("[提取文本摘要] 文本分析失败: %+v.\n", text.Error)
		task.Status = util.JobErrTextAnalysisFailed
		redis.UpdataOne(task)
		return
	}

	// 文本分析成功, 初始化需要存储的数据, 存入数据库
	abstext := abstext.NewAbsText(task.URL, task.KeyWords)
	abstext.Text = text.AText
	abstext.Abstract = text.TAbstract
	redis.InsertOne(abstext)
	mongodb.InsertOne(abstext)

	task.TextHash = abstext.Hash
	redis.UpdataOne(task)
	setAbstractFlag(task, util.JobTextAbstractExtractionDone)
}
