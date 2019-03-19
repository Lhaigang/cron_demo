package task

type WorkManager interface {
	CreateWork() Work
	Finish(task Task)
}

type Work interface {
	Init(task Task)
	RunWorker(task Task)
	Close(task Task)
}
type Task interface {
}
