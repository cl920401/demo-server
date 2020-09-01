package feishu

import (
	"demo-server/lib/config"
	"demo-server/lib/worker"
	"fmt"
	"net/http"
	"time"
)

//默认最大支持1024个待发送任务
var jobLength = 1024
var jobQueue = make(chan worker.Job, jobLength)

//2个发送线程
var dispatcher = worker.NewDispatcher(2)

func init() {
	dispatcher.Run(jobQueue)
}

func Send(title, text string) {
	//同步发送
	if config.Get("feishu.sync").Bool(true) {
		send(title, text)
		return
	}
	//TODO
	//当发送对列已满时，考虑到飞书发送限流等问题，新增发送任务直接丢弃
	if len(jobQueue) == jobLength {
		return
	}
	//异步发送
	job := worker.Job{
		Data: struct {
			Title string
			Text  string
		}{Title: title, Text: text},
		Proc: func(i interface{}) {
			data := i.(struct {
				Title string
				Text  string
			})
			send(data.Title, data.Text)
		},
	}

	select {
	case jobQueue <- job:
	default:
	}
}

func send(title, text string) {

	enable := config.Get("feishu.enable").Bool(false)
	if !enable {
		return
	}

	name := config.Get("service.name").String("")
	env := config.Get("service.env").String("")
	if name != "" && env != "" {
		title = fmt.Sprintf("【%s - %s】%s", name, env, title)
	}
	var data = struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	}{
		Title: title,
		Text:  text,
	}
	req := Request{TimeOut: 1 * time.Second}
	_, _ = req.Request(config.Get("feishu.webhook").String(""), http.MethodPost, data, nil)
}
