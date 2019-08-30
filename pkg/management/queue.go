package management

import (
	"errors"

	"github.com/mjosc/sqs/pkg/common"
)

func NewQueueManager(client common.SQSClient) common.SQSQueueManager {
	return &QueueManager{
		Client: client,
	}
}

type QueueManager struct {
	Client common.SQSClient
}

func (m *QueueManager) CreateQueue() error {
	return errors.New("not implemented")
}

func (m *QueueManager) DeleteQueue() error {
	return errors.New("not implemented")
}

func (m *QueueManager) ListQueues() ([]*string, error) {
	output, err := m.Client.ListQueues(nil)
	if err != nil {
		return nil, err // TODO: Wrap error
	}
	return output, nil
}
