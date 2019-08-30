package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/mjosc/sqs/pkg/client"
	"github.com/mjosc/sqs/pkg/components"
)

var queueURL = "some_url_here"

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		log.Fatalf("error initializing aws session: %v", err)
	}

	client := client.New(sess)
	producer := components.NewProducer(1, client, queueURL)
	consumer := components.NewConsumer(1, client, queueURL)

	if err = producer.Produce("Hello, World!"); err != nil {
		log.Fatalf("produce error: %v", err)
	}

	if err = consumer.Consume(); err != nil {
		log.Fatalf("consume error: %v", err)
	}
}
