package server

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"swc/logger"
	"swc/mongodb"
	"swc/mongodb/abstext"
	"swc/mongodb/absvideo"
	"swc/mongodb/job"
	"swc/mongodb/resource"
	"swc/util"
)

// textAbstract 用于存储文本分析的结果
type textAbstract struct {
	AText     string `json:"AText"`
	TAbstract string `json:"TAbstract"`
	Error     string `json:"Error"`
}

// textAnalysis 文本分析
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
	logger.Info.Println(python)
	go python.Call(job, textHandle)
}

// textHandle 文本分析的回调
func textHandle(job *job.Job, result []string) {
	// 获取资源信息
	r, err := resource.GetByKey(job.URL)
	if err != nil {
		// 获取资源出错
		job.SetStatus(util.JobErrFailedToFindResource)
		go JobSchedule(job)
		return
	}

	if len(result) != 1 {
		job.SetStatus(util.JobErrTextAnalysisFailed)
		go JobSchedule(job)
		return
	}

	var text textAbstract
	err = json.Unmarshal([]byte(result[0]), &text)
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
	r.SetAbsText(hash)
	job.SetAbsText(hash)
	mongodb.InsertOne(abstext)
}

// 从 url 和 keywords 中获取哈希
func getAbsHash(url string, keyWords []string) string {
	var str [12]byte
	hash := sha1.New()
	textStr := url + strings.Join(keyWords, "")
	hash.Write([]byte(textStr))
	copy(str[:], hash.Sum([]byte(""))[0:12])

	return fmt.Sprintf("%v", str)
}

// videoAbstract 用于存储视频摘要的结果
type videoAbstract struct {
	VAbstract []string `json:"VAbstract"`
	Error     string   `json:"Error"`
}

// videoAnalysis 视频分析
func videoAnalysis(job *job.Job) {
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
		PackagePath: filepath.Join(util.WorkSpace, "video_analysis"),
		FileName:    "api",
		MethodName:  "generate_abstract_from_video",
		Args: []string{
			SetArg(filepath.Join(r.Location, r.VideoPath)),
			SetArg(r.Location),
		},
	}

	// 文本分析
	logger.Info.Println(python)
	go python.Call(job, textHandle)
}

// videoHandle 视频分析的回调
func videoHandle(job *job.Job, result []string) {
	// 获取资源信息
	r, err := resource.GetByKey(job.URL)
	if err != nil {
		// 获取资源出错
		job.SetStatus(util.JobErrFailedToFindResource)
		go JobSchedule(job)
		return
	}

	if len(result) != 1 {
		job.SetStatus(util.JobErrTextAnalysisFailed)
		go JobSchedule(job)
		return
	}

	var videoPath videoAbstract
	err = json.Unmarshal([]byte(result[0]), &videoPath)
	if err != nil {
		job.SetStatus(util.JobErrTextAnalysisReadJSONFailed)
		go JobSchedule(job)
		return
	}
	hash := getAbsHash(job.URL, job.KeyWords)
	absvideo := absvideo.AbsVideo{
		URL:      job.URL,
		Abstract: videoPath.VAbstract,
	}
	r.SetAbsText(hash)
	job.SetAbsText(hash)
	mongodb.InsertOne(absvideo)
}
