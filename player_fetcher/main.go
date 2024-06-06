package main

import (
	"context"
	"log"
	"os"

	"github.com/ficontini/euro2024/player_fetcher/service"
	"github.com/joho/godotenv"
)

const (
	api_key_env     = "API_KEY"
	api_host_env    = "API_HOST"
	ws_endpoint_env = "WS_ENDPOINT"
)

func main() {
	var (
		fetcher = service.NewAPIFetcher(
			os.Getenv(api_host_env),
			os.Getenv(api_key_env))
		processor = service.NewAPIProcessor()
		svc       = service.New(fetcher, processor)
		endpoint  = os.Getenv(ws_endpoint_env)
	)
	client, err := NewWebSocketClient(endpoint, svc)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := client.SendMessage(ctx); err != nil {
		log.Fatal(err)
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
