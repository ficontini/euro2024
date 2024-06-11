package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ficontini/euro2024/matchfetcher/service"
	"github.com/ficontini/euro2024/matchfetcher/source/openliga"
	"github.com/joho/godotenv"
)

const (
	ws_endpoint_env = "WS_ENDPOINT"
	api_addr_env    = "OPENLIGA_ADDR"
)

var interval = 1 * time.Hour

func main() {
	var (
		fetcher   = openliga.NewAPIFetcher(os.Getenv(api_addr_env))
		processor = openliga.NewApiProcessor()
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
