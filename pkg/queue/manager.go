package queue

import (
	"errors"

	"github.com/mjosc/sqs/pkg/common"
)

func NewManager(client common.SQSClient) common.SQSQueueManager {
	return &Manager{
		Client: client,
	}
}

type Manager struct {
	Client common.SQSClient
}

func (m *Manager) CreateQueue() error {
	return errors.New("not implemented")
}

func (m *Manager) DeleteQueue() error {
	return errors.New("not implemented")
}

func (m *Manager) ListQueues() ([]*string, error) {
	output, err := m.Client.ListQueues(nil)
	if err != nil {
		return nil, err // TODO: Wrap error
	}
	return output, nil
}
