package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const (
	ws_endpoint_env = "WS_ENDPOINT"
	api_key_env     = "API_KEY"
	api_host_env    = "API_HOST"
	path            = "v3/fixtures?league=4&season=2024"
)

var interval = 15 * time.Minute

func main() {
	var (
		fetcher   = NewAPIFetcher(os.Getenv(api_host_env), os.Getenv(api_key_env), path)
		processor = NewApiProcessor()
		svc       = New(fetcher, processor)
		endpoint  = os.Getenv(ws_endpoint_env)
	)
	client, err := NewWebSocketClient(endpoint, svc)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		if err := client.SendMessage(ctx); err != nil {
			log.Fatal(err)
		}
		time.Sleep(interval)
	}

}
func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
