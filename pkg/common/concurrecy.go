package common

type ThreadPool interface {
	Execute(Task)
	Close()
}

type Task interface {
	Run()
}
