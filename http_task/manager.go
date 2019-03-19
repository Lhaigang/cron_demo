package http_task

import (
	"cron_demo/task"
	"cron_demo/work"
	"io"
	"os"
)

// Run makes all the requests, prints the summary. It blocks until
// all work is done.
func NewManager(t task.Task) task.WorkManager {
	// append hey's user agent
	this := new(Manager)
	//this.results = make(chan *work.Result, t.N)
	return this
}

type Manager struct {
	results chan *work.Result
	// Writer is where results will be written. If nil, results are written to stdout.
	Writer io.Writer
}

func (this *Manager) writer() io.Writer {
	if this.Writer == nil {
		return os.Stdout
	}
	return this.Writer
}
func (this *Manager) Finish(task task.Task) {
	close(this.results)
	//total := time.Now().Sub(task.Start)
	//work.NewReport(this.writer(), task.N, this.results, "", total).Finalize()
}
func (this *Manager) CreateWork() task.Work {
	w := new(Work)
	w.H2 = false
	w.Timeout = 20

	w.DisableCompression = false
	w.DisableKeepAlives = false
	w.DisableRedirects = false
	w.manager = this
	return w
}
