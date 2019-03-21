package http_task

import (
	"cron_demo/http_request"
	"cron_demo/task"
	"cron_demo/work"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

/**
Work 代表一个协程内具体执行任务工作者
*/
type Work struct {
	work.HttpWork
	manager  *Manager
	QPS      int
	closeSig bool
	num      int
}


func (this *Work) Init(t task.Task) {
	this.QPS = 10
	this.num = 0
	this.closeSig = false
}

/**

 */
func (this *Work) RunWorker(t task.Task) {
	for !this.closeSig {
		var throttle <-chan time.Time
		if this.QPS > 0 {
			throttle = time.Tick(time.Duration(1e6/(this.QPS)) * time.Microsecond)
		}

		if this.QPS > 0 {
			<-throttle
		}
		this.worker(t)
	}
}

//不同的请求

func (this *Work) worker(t task.Task) {
	this.num++
	//发送模板消息
	http_request.SendTemplate()
}
func (this *Work) Close(t task.Task) {
	this.closeSig = true
	fmt.Println(fmt.Sprintf("num : %d", this.num))
}








