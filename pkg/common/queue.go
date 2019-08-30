package common

type SQSQueueManager interface {
	CreateQueue()
	DeleteQueue()
}
