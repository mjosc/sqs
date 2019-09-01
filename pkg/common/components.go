package common

import "log"

type Producer interface {
	Produce(messageBody string) error
}

type Consumer interface {
	Consume() error
}

func NewConsumeTask(consumer Consumer) *ConsumeTask {
	return &ConsumeTask{
		Consumer: consumer,
	}
}

type ConsumeTask struct {
	Consumer Consumer
}

func (t *ConsumeTask) Run() {
	if err := t.Consumer.Consume(); err != nil {
		log.Fatalln(err) // TODO
	}
}
