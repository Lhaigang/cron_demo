package task

import (
	"cron_demo/utils"
	"sync"
	"time"
)

type LoopTask struct {
	// C is the concurrency level, the number of concurrent workers to run.
	C int

	Start time.Time

	q *utils.Queue

	wg sync.WaitGroup
}

// Run makes all the requests, prints the summary. It blocks until
// all work is done.
func (b *LoopTask) Run(manager WorkManager) {
	b.Start = time.Now()
	b.q = utils.NewQueue()
	b.runWorkers(manager)
}

func (b *LoopTask) Stop() {
	end := false
	for !end {
		worker := b.q.Pop()
		if worker != nil {
			worker.(Work).Close(b)
		} else {
			end = true
		}
	}

}
func (b *LoopTask) Wait() {
	b.wg.Wait()
}
func (b *LoopTask) runWorkers(manager WorkManager) {
	b.wg.Add(b.C)
	for i := 0; i < b.C; i++ {
		task := manager.CreateWork()
		b.q.Push(task)
		go func() {
			task.Init(b)
			task.RunWorker(b)
			b.wg.Done()
		}()
	}

}
