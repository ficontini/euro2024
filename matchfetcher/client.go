package main

import (
	"context"
	"encoding/json"

	"github.com/ficontini/euro2024/matchfetcher/service"
	"github.com/ficontini/euro2024/types"
	"github.com/gorilla/websocket"
)

const default_action = "sendMessage"

type WebSocketClient struct {
	service service.Service
	conn    *websocket.Conn
}

func NewWebSocketClient(endpoint string, service service.Service) (*WebSocketClient, error) {
	conn, _, err := websocket.DefaultDialer.Dial(endpoint, nil)
	if err != nil {
		return nil, err
	}
	return &WebSocketClient{
		service: service,
		conn:    conn,
	}, nil
}

func (c *WebSocketClient) SendMessage(ctx context.Context) error {
	matches, err := c.service.FetchMatches(ctx)
	if err != nil {
		return err
	}
	messageJSON, err := json.Marshal(NewMessage(matches))
	if err != nil {
		return err
	}
	return c.conn.WriteMessage(websocket.TextMessage, messageJSON)

}
func (c *WebSocketClient) Close() error {
	return c.conn.Close()
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
