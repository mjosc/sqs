package client

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/mjosc/sqs/pkg/common"
)

func New(sess *session.Session) common.SQSClient {
	return &Client{
		Service: sqs.New(sess),
	}
}

type Client struct {
	Service *sqs.SQS
}

func (c *Client) SendMessage(input *sqs.SendMessageInput) (string, error) {
	result, err := c.Service.SendMessage(input)
	if err != nil {
		return "", fmt.Errorf("message send error: %v", err)
	}
	messageID := result.MessageId
	if messageID == nil {
		return "", errors.New("message send error: unexpected nil id")
	}
	return *messageID, nil
}

func (c *Client) ReceiveMessage(input *sqs.ReceiveMessageInput) ([]*sqs.Message, error) {
	result, err := c.Service.ReceiveMessage(input)
	if err != nil {
		return nil, fmt.Errorf("message receive error: %v", err)
	}
	return result.Messages, nil
}

func (c *Client) DeleteMessage() error {
	return nil
}
