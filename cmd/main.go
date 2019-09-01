package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/mjosc/sqs/pkg/components"
	"github.com/mjosc/sqs/pkg/concurrency"

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

	rand.Seed(time.Now().UnixNano())

	c := components.NewConsumer(1, client, *url)
	consumer := components.NewMultiThreadedConsumer(c, 10)

	go func() {
		producer := components.NewProducer(1, client, *url)
		pool := concurrency.NewPool(4)

		for i := 0; i < 100; i++ {
			msg := fmt.Sprintf("Hello, SQS! [%d]", i)
			task := components.NewProducerTask(producer, msg)
			pool.Execute(task)
		}

		pool.Close()
	}()

	if err = consumer.Consume(); err != nil {
		log.Fatalln(err)
	}
}
