package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

const (
	endpoint_env_var = "endpoint"
	queueURL_env_var = "queue_url"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	var (
		client = apigatewaymanagementapi.NewFromConfig(cfg, func(o *apigatewaymanagementapi.Options) {
			o.BaseEndpoint = aws.String(os.Getenv(endpoint_env_var))
		})
		sqsClient = sqs.NewFromConfig(cfg)
		producer  = NewSQSProducer(os.Getenv(queueURL_env_var), sqsClient)
		handler   = NewHandler(client, producer)
	)

	lambda.Start(handler.HandleMessage)
}
