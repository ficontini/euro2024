package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ficontini/euro2024/types"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

const (
	wsEndpoint   = "ws://127.0.0.1:3000/ws"
	origin       = "http://localhost/"
	api_key_env  = "API_KEY"
	api_host_env = "API_HOST"
	path         = "v3/fixtures?league=4&season=2024"
)

var interval = 5 * time.Minute

func main() {
	var (
		fetcher   = NewAPIFetcher(os.Getenv(api_host_env), os.Getenv(api_key_env), path)
		processor = NewApiProcessor()
		svc       = New(fetcher, processor)
	)
	//_ = svc
	matches, err := svc.FetchLiveMatches(context.Background())
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
func fakeMatches() []*types.Match {
	var matches []*types.Match
	for i := 0; i < 40; i++ {
		matches = append(matches,
			types.NewMatch(
				time.Now(),
				types.NewLocation("Munich", "Germany"),
				"Italy",
				"Albania",
				types.NS,
				types.NewResult(0, 0),
			))
	}
	return matches
}
