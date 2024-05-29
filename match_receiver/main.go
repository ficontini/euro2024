package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/joho/godotenv"
)

const queue_url_env_var = "QUEUE_URL"

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	var (
		client   = sqs.NewFromConfig(cfg)
		producer = New(os.Getenv(queue_url_env_var), client)
		server   = NewServer(producer)
	)

	log.Fatal(server.Start())
}
func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
