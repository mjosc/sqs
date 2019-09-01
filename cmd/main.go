package main

import (
	"fmt"
	"log"

	"github.com/mjosc/sqs/pkg/components"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/mjosc/sqs/pkg/client"
	"github.com/mjosc/sqs/pkg/queue"
)

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		log.Fatalf("error initializing aws session: %v\n", err)
	}

	client := client.New(sess)
	manager := queue.NewManager(client)

	queues, err := manager.ListQueues()
	if err != nil {
		log.Fatalf("error retrieving queue list: %v\n", err)
	}

	if len(queues) < 1 {
		log.Fatalln("no queues exist")
	}
	url := queues[0]
	if url == nil {
		log.Fatalln("unexpected nil queue url")
	}

	producer := components.NewProducer(1, client, *url)

	c := components.NewConsumer(1, client, *url)
	consumer := components.NewMultiThreadedConsumer(c)

	for i := 0; i < 5; i++ {
		message := fmt.Sprintf("Hello, SQS! [%d]", i)
		producer.Produce(message)
	}

	fmt.Println("all messages have been added to the queue")

	// Consume is a blocking method. It will forever poll for new messages.
	if err = consumer.Consume(); err != nil {
		log.Fatalln(err)
	}
}
