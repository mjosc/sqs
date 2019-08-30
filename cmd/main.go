package main

import (
	"log"
	"strconv"

	"github.com/mjosc/sqs/pkg/common"

	"github.com/mjosc/sqs/pkg/components"
	"github.com/mjosc/sqs/pkg/management"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/mjosc/sqs/pkg/client"
)

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		log.Fatalf("error initializing aws session: %v\n", err)
	}

	client := client.New(sess)
	manager := management.NewQueueManager(client)

	queues, err := manager.ListQueues()
	if err != nil {
		log.Fatalf("error retrieving queue list: %v\n", err)
	}

	producers := make([]common.Producer, 0, 10)
	consumers := make([]common.Consumer, 0, 10)

	for i, q := range queues {
		if q == nil {
			continue
		}
		producers = append(producers, components.NewProducer(i, client, *q))
		consumers = append(consumers, components.NewConsumer(i, client, *q))
	}

	for i := 0; i < 10; i++ {
		if err = producers[0].Produce(strconv.Itoa(i)); err != nil {
			log.Fatalln(err)
		}
		if err = consumers[0].Consume(); err != nil {
			log.Fatalln(err)
		}
	}
}
