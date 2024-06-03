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
	matchendpoint "github.com/ficontini/euro2024/matchservice/pkg/endpoint"
	"github.com/ficontini/euro2024/matchservice/pkg/service"
	"github.com/ficontini/euro2024/matchservice/pkg/transport"
	"github.com/ficontini/euro2024/matchservice/proto"
	"github.com/ficontini/euro2024/matchservice/queue"
	"github.com/ficontini/euro2024/matchservice/store"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

const (
	queue_url_env     = "QUEUE_URL"
	http_listener_env = "HTTP_LISTENER"
	grpc_listener_env = "GRPC_LISTENER"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	var (
		client      = sqs.NewFromConfig(cfg)
		store       = store.NewInMemoryStore()
		svc         = service.New(store)
		queueSvc    = queue.New(store)
		consumer    = queue.NewSQSConsumer(client, os.Getenv(queue_url_env), queueSvc)
		endpoints   = matchendpoint.New(svc)
		httpHandler = transport.NewHTTPHandler(endpoints)
		grpcHandler = transport.NewGRPCServer(endpoints)
	)

	go func() {
		httpListener, err := net.Listen("tcp", os.Getenv(http_listener_env))
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		fmt.Println("starting server on:", httpListener.Addr())
		if err := http.Serve(httpListener, httpHandler); err != nil {
			panic(err)
		}
	}()

	go func() {
		ln, err := net.Listen("tcp", os.Getenv(grpc_listener_env))
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		server := grpc.NewServer(grpc.EmptyServerOption{})
		proto.RegisterMatchesServer(server, grpcHandler)
		if err := server.Serve(ln); err != nil {
			log.Fatal(err)
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
