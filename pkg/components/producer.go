package components

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/mjosc/sqs/pkg/common"
)

func NewProducerTask(producer common.Producer, message string) common.Task {
	return &ProducerTask{
		Producer: producer,
		Message:  message,
	}
}

type ProducerTask struct {
	Producer common.Producer
	Message  string
}

func (p *ProducerTask) Run() {
	r := rand.Intn(3)
	time.Sleep(time.Duration(r) * time.Second)
	p.Producer.Produce(p.Message)
}

func NewProducer(id int, client common.SQSClient, url string) common.Producer {
	return &Producer{
		ID:       id,
		Client:   client,
		QueueURL: url,
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
	fmt.Printf("sent message: %v\n", msg)
	return nil
}
