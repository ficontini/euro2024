package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/ficontini/euro2024/playerstorer/store"
	"github.com/joho/godotenv"
)

const queue_url_env = "PLAYER_QUEUE_URL"

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	var (
		sqsClient      = sqs.NewFromConfig(cfg)
		dynamodbClient = dynamodb.NewFromConfig(cfg)
		store          = store.NewDynamoDBStore(dynamodbClient)
		svc            = New(store)
		consumer       = NewSQSConsumer(sqsClient, os.Getenv(queue_url_env), svc)
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer.Start(ctx)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	consumer.Stop(ctx)
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
