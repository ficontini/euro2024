package main

import (
	"log"
	"os"

	"github.com/ficontini/euro2024/gateway/api"
	"github.com/ficontini/euro2024/matchservice/pkg/transport"
	playertransport "github.com/ficontini/euro2024/playerservice/transport"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	matchGrpcAddrEnvVar  = "GRPC_LISTENER"
	playerGrpcAddrEnvVar = "PLAYER_GRPC_LISTENER"
	listenAddrEnVar      = "GATEWAY_ADDR"
)

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {
	matchServiceConn, err := newGRPCConnection(getEnv(matchGrpcAddrEnvVar))
	if err != nil {
		log.Fatal(err)
	}
	defer matchServiceConn.Close()

	playerServiceConn, err := newGRPCConnection(getEnv(playerGrpcAddrEnvVar))
	if err != nil {
		log.Fatal(err)
	}
	defer playerServiceConn.Close()

	var (
		listenAddr    = getEnv(listenAddrEnVar)
		matchService  = transport.NewGRPCClient(matchServiceConn)
		playerService = playertransport.NewGRPCClient(playerServiceConn)
		app           = fiber.New(config)
		apiv1         = app.Group("/api/v1")
		matches       = apiv1.Group("/matches")
		matchHandler  = api.NewMatchHandler(matchService)
		playerHandler = api.NewPlayerHandler(playerService)
	)

	matches.Get("/upcoming", matchHandler.HandleGetUpcomingMatches)
	matches.Get("/live", matchHandler.HandleGetLiveMatches)
	matches.Get("/:team", matchHandler.HandleGetMatchesByTeam)
	matches.Get("/:team/players", playerHandler.HandleGetPlayersByTeam)
	log.Fatal(app.Listen(listenAddr))
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("env var %s not set", key)
	}
	return value
}
func newGRPCConnection(addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
