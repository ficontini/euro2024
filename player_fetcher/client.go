package main

import (
	"context"
	"encoding/json"

	"github.com/ficontini/euro2024/player_fetcher/service"
	"github.com/ficontini/euro2024/types"
	"github.com/gorilla/websocket"
)

const default_action = "sendPlayers"

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
	players, err := c.service.FetchData()
	if err != nil {
		return err
	}
	messageJSON, err := json.Marshal(NewMessage(players))
	if err != nil {
		return err
	}
	return c.conn.WriteMessage(websocket.TextMessage, messageJSON)

}
func (c *WebSocketClient) Close() error {
	return c.conn.Close()
}

type Message struct {
	Action  string          `json:"action"`
	Players []*types.Player `json:"players"`
}

func NewMessage(players []*types.Player) Message {
	return Message{
		Action:  default_action,
		Players: players,
	}
}
