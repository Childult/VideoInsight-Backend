package server

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"swc/mongodb"
	"swc/mongodb/abstext"
	"swc/mongodb/job"
	"swc/mongodb/resource"
	"swc/util"
)

type textAbstract struct {
	AText     string `json:"AText"`
	TAbstract string `json:"TAbstract"`
	Error     string `json:"Error"`
}

func textAnalysis(job *job.Job) {
	// 获取资源信息
	r, err := resource.GetByKey(job.URL)
	if err != nil {
		// 获取资源出错
		job.SetStatus(util.JobErrFailedToFindResource)
		go JobSchedule(job)
		return
	}

	// 构建文本分析对象
	python := PyWorker{
		PackagePath: filepath.Join(util.WorkSpace, "text_analysis"),
		FileName:    "api",
		MethodName:  "generate_abstract_from_audio",
		Args: []string{
			SetArg(filepath.Join(r.Location, r.AudioPath)),
		},
	}

	// 文本分析
	fmt.Println(python)
	// go python.Call(job, textHandle)
}

func textHandle(job *job.Job, result []string) {
	if len(result) != 1 {
		job.SetStatus(util.JobErrTextAnalysisFailed)
		go JobSchedule(job)
		return
	}

	var text textAbstract
	err := json.Unmarshal([]byte(result[0]), &text)
	if err != nil {
		job.SetStatus(util.JobErrTextAnalysisReadJSONFailed)
		go JobSchedule(job)
		return
	}
	hash := getAbsHash(job.URL, job.KeyWords)
	abstext := abstext.AbsText{
		Hash:     hash,
		URL:      job.URL,
		KeyWords: job.KeyWords,
		Text:     text.AText,
		Abstract: text.TAbstract,
	}
	job.SetAbsText(hash)
	mongodb.InsertOne(abstext)
}

func getAbsHash(url string, keyWords []string) string {
	var str [12]byte
	hash := sha1.New()
	textStr := url + strings.Join(keyWords, "")
	hash.Write([]byte(textStr))
	copy(str[:], hash.Sum([]byte(""))[0:12])

	return fmt.Sprintf("%v", str)
}
