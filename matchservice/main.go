package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	matchendpoint "github.com/ficontini/euro2024/matchservice/endpoint"
	"github.com/ficontini/euro2024/matchservice/service"
	"github.com/ficontini/euro2024/matchservice/transport"
	"github.com/joho/godotenv"
)

const queue_url_env = "QUEUE_URL"

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	var (
		client      = sqs.NewFromConfig(cfg)
		svc         = service.New()
		consumer    = transport.NewSQSConsumer(client, os.Getenv(queue_url_env), svc)
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
