package common

import (
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQSClient interface {
	ListQueues(input *sqs.ListQueuesInput) ([]*string, error)
	SendMessage(input *sqs.SendMessageInput) (string, error)
	ReceiveMessage(input *sqs.ReceiveMessageInput) ([]*sqs.Message, error)
	DeleteMessage(input *sqs.DeleteMessageInput) error
}
