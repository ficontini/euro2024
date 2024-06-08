package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	playerendpoint "github.com/ficontini/euro2024/playerservice/endpoint"
	"github.com/ficontini/euro2024/playerservice/proto"
	"github.com/ficontini/euro2024/playerservice/service"
	"github.com/ficontini/euro2024/playerservice/store"
	"github.com/ficontini/euro2024/playerservice/transport"
)

const playerGrpcListenerEnv = "PLAYER_GRPC_LISTENER"

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	var (
		client      = dynamodb.NewFromConfig(cfg)
		store       = store.NewDynamoDBStore(client)
		svc         = service.New(store)
		endpoint    = playerendpoint.New(svc)
		grpcHandler = transport.NewGRPCServer(endpoint)
	)

	ln, err := net.Listen("tcp", os.Getenv(playerGrpcListenerEnv))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	server := grpc.NewServer(grpc.EmptyServerOption{})
	proto.RegisterPlayersServer(server, grpcHandler)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("shutting down the server")
		server.GracefulStop()
		ln.Close()
		log.Println("server gracefully stopped.")
	}()

	if err := server.Serve(ln); err != nil {
		log.Fatal(err)
	}

}
func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
