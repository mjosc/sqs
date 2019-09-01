package common

type Producer interface {
	Produce(messageBody string) error
}

type Consumer interface {
	Consume() error
}
