package components

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/mjosc/sqs/pkg/common"
)

func NewProducer(id int, client common.SQSClient, url string) common.Producer {
	return &Producer{
		ID:       id,
		Client:   client,
		QueueURL: url, // TODO: Make these private fields so as to be thread safe
	}
}

type Producer struct {
	ID       int
	Client   common.SQSClient
	QueueURL string
}

func (p *Producer) Produce(msg string) error {
	stringID := strconv.Itoa(p.ID)

	input := sqs.SendMessageInput{
		QueueUrl: aws.String(p.QueueURL),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"ProducerID": &sqs.MessageAttributeValue{
				DataType:    aws.String("Number"),
				StringValue: aws.String(stringID),
			},
		},
		MessageBody: aws.String(msg),
	}

	_, err := p.Client.SendMessage(&input)
	if err != nil {
		return err // TODO: Wrap error
	}
	return nil
}
