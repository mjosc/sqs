package common

import (
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQSClient interface {
	SendMessage(input *sqs.SendMessageInput) (string, error)
	ReceiveMessage(input *sqs.ReceiveMessageInput) ([]*sqs.Message, error)
	DeleteMessage() error
}
