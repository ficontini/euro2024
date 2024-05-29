package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

const (
	wsEndpoint   = "ws://127.0.0.1:3000/ws"
	api_key_env  = "API_KEY"
	api_host_env = "API_HOST"
	path         = "v3/fixtures?league=4&season=2024"
)

var interval = 15 * time.Minute

func main() {
	var (
		fetcher   = NewAPIFetcher(os.Getenv(api_host_env), os.Getenv(api_key_env), path)
		processor = NewApiProcessor()
		svc       = New(fetcher, processor)
	)

	matches, err := svc.FetchMatches(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		if err := conn.WriteJSON(matches); err != nil {
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
