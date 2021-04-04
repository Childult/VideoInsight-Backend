package task_builder

import (
	"swc/data/abstext"
	"swc/data/resource"
	"swc/data/task"
	"swc/dbs/mongodb"
	"swc/dbs/redis"
	"swc/logger"
	"swc/server/abstext_builder"
	"swc/server/absvideo_builder"
	"swc/server/resource_builder"
	"swc/util"
	"sync"
)

// TaskArgs 任务传递时的参数
type TaskArgs struct {
	t *task.Task         // 任务
	r *resource.Resource // 资源
}

// TaskScheduler 任务调度器
type TaskScheduler struct {
	mu    sync.Mutex     // 任务锁, 保证任务的创建唯一
	tasks chan *TaskArgs // 任务管道, 保证任务调度串行执行
}

var ts TaskScheduler
var onceTS sync.Once

// AddTask 新增任务
// url: 要获取的资源链接
// keyWords: 关键词
func AddTask(url string, keyWords []string) {
	onceTS.Do(func() {
		ts.tasks = make(chan *TaskArgs)
		go ts.Scheduler(ts.tasks)
	})

	// 构建参数
	t := task.NewTask(url, keyWords)
	r := &resource.Resource{URL: url}
	arg := &TaskArgs{t: t, r: r}
	ts.tasks <- arg
	logger.Debug.Println("[任务调度] 收到任务", t)
}

// Scheduler 核心任务的调度器, 串行执行
func (tb *TaskScheduler) Scheduler(args chan *TaskArgs) {
	// 等待任务
	for arg := range args {
		t := arg.t
		switch t.Status {
		case util.TaskStart:
			go tb.createTask(arg)
		case util.TaskCreated:
			go tb.RetrieveResource(arg)
		case util.TaskDownloadDone:
			ch := make(chan int32)
			go tb.textAnalysis(ch, arg)
			go tb.videoAnalysis(ch, arg)
			go tb.combine(ch, arg)
		case util.TaskCompleted:
			logger.Info.Println("任务完成")
		default:
			logger.Error.Println("发生错误:", t)
		}
		// r, err := ts.GetResource(t.URL)
	}
}

// createTask 判断任务是否存在, 不存在则创建任务
func (tb *TaskScheduler) createTask(arg *TaskArgs) {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	if redis.Exists(arg.t) || mongodb.Exists(arg.t) { // 任务已经存在, 直接忽略
		return
	}
	// 不存在则正式开始任务
	arg.t.Status = util.TaskCreated
	redis.InsertOne(arg.t)
	tb.tasks <- arg

	// 关键词不为空时, 会发起两个任务
	if arg.t.KeyWords != nil {
		newTask := task.NewTask(arg.t.URL, nil)
		newRS := &resource.Resource{URL: arg.t.URL}
		newArg := &TaskArgs{t: newTask, r: newRS}
		tb.tasks <- newArg
	}
}

// RetrieveResource 获取资源
func (tb *TaskScheduler) RetrieveResource(arg *TaskArgs) {
	err := resource_builder.RequestResource(arg.t.URL)
	if err == nil {
		arg.t.Status = util.TaskDownloadDone
		redis.FindOne(arg.r)
	} else {
		arg.t.Status = util.TaskErrRetrieveFail
	}
	redis.UpdataOne(arg.t)
	tb.tasks <- arg
}

// textAnalysis 调用文本分析, 成功时往 ch 发 1, 否则发 0
func (tb *TaskScheduler) textAnalysis(ch chan int32, arg *TaskArgs) {
	logger.Debug.Println("[任务调度] 文本分析开始")
	err := abstext_builder.RequestTextAnalysis(arg.t.URL, arg.t.KeyWords, arg.r.VideoPath)
	if err == nil {
		at := abstext.NewAbsText(arg.t.URL, arg.t.KeyWords)
		redis.FindOne(at)
		arg.t.TextHash = at.Hash
		redis.UpdataOne(arg.t)
		ch <- 1
	} else {
		ch <- 0
	}
	logger.Debug.Println("[任务调度] 文本分析结束", err)
}

// videoAnalysis 调用视频分析, 成功时往 ch 发 2, 否则发 0
func (tb *TaskScheduler) videoAnalysis(ch chan int32, arg *TaskArgs) {
	logger.Debug.Println("[任务调度] 视频分析开始")
	err := absvideo_builder.RequestVideoAnalysis(arg.t.URL, arg.r.VideoPath, arg.r.Location)
	if err == nil {
		ch <- 2
	} else {
		ch <- 0
	}
	logger.Debug.Println("[任务调度] 视频分析结束", err)
}

// combine 合并结果
func (tb *TaskScheduler) combine(ch chan int32, arg *TaskArgs) {
	logger.Debug.Println("[任务调度] 等待文本分析和视频分析结束")
	status1 := <-ch
	status2 := <-ch
	status := status1 | status2
	switch status {
	case 0: // 文本分析和视频分析都失败
		arg.t.Status = util.TaskErrTextAndVideoFail
	case 1: // 文本分析成功, 视频分析失败
		arg.t.Status = util.TaskErrVideoAnalysisFailed
	case 2: // 文本分析失败, 视频分析成功
		arg.t.Status = util.TaskErrTextAnalysisFailed
	case 3: // 文本分析和视频分析都成功
		arg.t.Status = util.TaskCompleted
		mongodb.InsertOne(arg.t)
	}
	redis.UpdataOne(arg.t)
	tb.tasks <- arg
}
