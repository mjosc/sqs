package components

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/mjosc/sqs/pkg/common"
)

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

	for _, msg := range output {

		messageID := msg.MessageId
		if messageID == nil {
			return errors.New("consumer error: message id not present")
		}
		producerAttribute, ok := msg.MessageAttributes["ProducerID"]
		if !ok || producerAttribute == nil {
			return errors.New("consumer error: producer id not present")
		}
		producerID := producerAttribute.StringValue
		if aws.StringValue(producerID) != strconv.Itoa(c.ID) {
			return fmt.Errorf("consumer error: incorrect producer id got %v but expected %v", *producerID, c.ID)
		}

		input := sqs.DeleteMessageInput{
			QueueUrl:      aws.String(c.QueueURL),
			ReceiptHandle: msg.ReceiptHandle,
		}
		if err := c.Client.DeleteMessage(&input); err != nil {
			return fmt.Errorf("consumer error: %v", err)
		}

		fmt.Printf("successfully processed message %v\n", *messageID)
	}

	return nil
}
