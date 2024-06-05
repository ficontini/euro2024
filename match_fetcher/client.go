package main

import (
	"context"
	"encoding/json"

	"github.com/ficontini/euro2024/match_fetcher/service"
	"github.com/ficontini/euro2024/types"
	"github.com/gorilla/websocket"
)

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
	messageJSON, err := json.Marshal(types.NewMessage(matches))
	if err != nil {
		return err
	}
	return c.conn.WriteMessage(websocket.TextMessage, messageJSON)

}
func (c *WebSocketClient) Close() error {
	return c.conn.Close()
}
