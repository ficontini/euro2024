package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	matchendpoint "github.com/ficontini/euro2024/matchservice/endpoint"
	"github.com/ficontini/euro2024/matchservice/service"
	"github.com/ficontini/euro2024/matchservice/transport"
	"github.com/joho/godotenv"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	var (
		client      = sqs.NewFromConfig(cfg)
		svc         = service.New()
		endpoints   = matchendpoint.New(svc)
		httpHandler = transport.NewHTTPHandler(endpoints)
	)
	go func() {
		httpListener, err := net.Listen("tcp", ":3003")
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		fmt.Println("starting server")
		if err := http.Serve(httpListener, httpHandler); err != nil {
			panic(err)
		}
	}()
	transport.NewSQSConsumer(client, os.Getenv("QUEUE_URL"), svc)
	select {}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
