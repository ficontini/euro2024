package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/ficontini/euro2024/types"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

const (
	ws_endpoint_env = "WS_ENDPOINT"
	api_key_env     = "API_KEY"
	api_host_env    = "API_HOST"
	default_action  = "sendMessage"
	path            = "v3/fixtures?league=4&season=2024"
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
	conn, _, err := websocket.DefaultDialer.Dial(os.Getenv(ws_endpoint_env), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	messageJSON, err := json.Marshal(NewMessage(matches))
	if err != nil {
		log.Fatal("marshalling err:", err)
	}
	for {
		if err := conn.WriteMessage(websocket.TextMessage, messageJSON); err != nil {
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

type Message struct {
	Action  string         `json:"action"`
	Matches []*types.Match `json:"matches"`
}

func NewMessage(matches []*types.Match) Message {
	return Message{
		Action:  default_action,
		Matches: matches,
	}
}
