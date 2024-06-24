package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/ficontini/euro2024/playerstorer/store"
	"github.com/ficontini/euro2024/util"
)

const queueUrlEnvVar = "PLAYER_QUEUE_URL"

var queueURL string

func main() {
	if err := Init(); err != nil {
		log.Fatal(err)
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	var (
		sqsClient      = sqs.NewFromConfig(cfg)
		dynamodbClient = dynamodb.NewFromConfig(cfg)
		store          = store.NewDynamoDBStore(dynamodbClient)
		svc            = New(store)
		consumer       = NewSQSConsumer(sqsClient, queueURL, svc)
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer.Start(ctx)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	consumer.Stop(ctx)
}

func Init() error {
	if err := util.Load(); err != nil {
		return err
	}
	queueURL = os.Getenv(os.Getenv(queueUrlEnvVar))
	if queueURL == "" {
		return fmt.Errorf("error loading the var: %s", queueUrlEnvVar)
	}
	return nil
}
