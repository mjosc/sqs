package components

import (
	"fmt"

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
		// process request here
		fmt.Printf("successfully processed request %v with content: %v\n", *msg.MessageId, *msg.Body)
		if err := c.Client.DeleteMessage(); err != nil {
			// do something
		}
	}

	return nil
}
