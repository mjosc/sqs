package common

type SQSQueueManager interface {
	CreateQueue() error
	DeleteQueue() error
	ListQueues() ([]*string, error)
}
