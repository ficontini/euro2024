package main

import (
	"log"
	"os"

	"github.com/ficontini/euro2024/gateway/api"
	"github.com/ficontini/euro2024/matchservice/pkg/transport"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpc_addr_env_var  = "GRPC_LISTENER"
	listen_addr_en_var = "GATEWAY_ADDR"
)

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {
	listenAddr := os.Getenv(listen_addr_en_var)
	grpcAddr := os.Getenv(grpc_addr_env_var)
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	var (
		svc     = transport.NewGRPCClient(conn)
		app     = fiber.New(config)
		apiv1   = app.Group("/api/v1")
		matches = apiv1.Group("/matches")
		handler = api.NewMatchHandler(svc)
	)
	matches.Get("/upcoming", handler.HandleGetUpcomingMatches)
	matches.Get("/live", handler.HandleGetLiveMatches)
	matches.Get("/:team", handler.HandleGetMatchesByTeam)
	log.Fatal(app.Listen(listenAddr))
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
