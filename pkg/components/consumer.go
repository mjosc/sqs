package components

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/mjosc/sqs/pkg/concurrency"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/mjosc/sqs/pkg/common"
)

func NewConsumerTask(consumer common.Consumer) common.Task {
	return &ConsumerTask{
		Consumer: consumer,
	}
}

type ConsumerTask struct {
	Consumer common.Consumer
}

func (t *ConsumerTask) Run() {
	if err := t.Consumer.Consume(); err != nil {
		log.Fatalln(err) // TODO
	}
}

func NewMultiThreadedConsumer(c common.Consumer, nThreads int) common.Consumer {
	return &MultiThreadedConsumer{
		Consumer:   c,
		ThreadPool: concurrency.NewPool(nThreads),
	}
}

type MultiThreadedConsumer struct {
	Consumer   common.Consumer
	ThreadPool common.ThreadPool
}

func (c *MultiThreadedConsumer) Consume() error {
	// continuously poll for messages
	for {
		task := NewConsumerTask(c.Consumer)
		c.ThreadPool.Execute(task)
	}
	return nil // TODO
}

func NewConsumer(id int, client common.SQSClient, url string) common.Consumer {
	return &Consumer{
		ID:       id,
		Client:   client,
		QueueURL: url,
	}
}

type Consumer struct {
	ID       int
	Client   common.SQSClient
	QueueURL string
}

func (c *Consumer) Consume() error {
	input := sqs.ReceiveMessageInput{
		QueueUrl: aws.String(c.QueueURL),
		AttributeNames: aws.StringSlice([]string{
			sqs.MessageSystemAttributeNameSentTimestamp,
		}),
		MessageAttributeNames: aws.StringSlice([]string{
			sqs.QueueAttributeNameAll,
		}),
	}

	output, err := c.Client.ReceiveMessage(&input)
	if err != nil {
		return err
	}

	if len(output) == 0 {
		fmt.Println("no messsages to consume")
		return nil
	}

	for _, msg := range output {
		messageID, err := c.processMessage(msg)
		if err != nil {
			return err
		}
		input := sqs.DeleteMessageInput{
			QueueUrl:      aws.String(c.QueueURL),
			ReceiptHandle: msg.ReceiptHandle,
		}
		if err := c.Client.DeleteMessage(&input); err != nil {
			return fmt.Errorf("consumer error: %v", err)
		}
		fmt.Printf("successfully processed message %v: %v\n", messageID, *msg.Body) // TODO
	}
	return nil
}

func (c *Consumer) processMessage(msg *sqs.Message) (string, error) {
	messageID := msg.MessageId
	if messageID == nil {
		return "", errors.New("consumer error: message id not present")
	}
	producerAttribute, ok := msg.MessageAttributes["ProducerID"]
	if !ok || producerAttribute == nil {
		return "", errors.New("consumer error: producer id not present")
	}
	producerID := producerAttribute.StringValue
	if aws.StringValue(producerID) != strconv.Itoa(c.ID) {
		return "", fmt.Errorf("consumer error: incorrect producer id got %v but expected %v", *producerID, c.ID)
	}
	return *messageID, nil
}
