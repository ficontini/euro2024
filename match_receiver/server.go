package main

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/ficontini/euro2024/types"
	"github.com/gorilla/websocket"
)

const buffer_size = 1028

type Server struct {
	conn     *websocket.Conn
	producer Producer
}

func NewServer(producer Producer) *Server {
	return &Server{
		producer: producer,
	}
}
func (s *Server) Start() error {
	http.HandleFunc("/ws", s.handleWS)
	return http.ListenAndServe(":3000", nil)
}
func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  buffer_size,
		WriteBufferSize: buffer_size,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	s.conn = conn
	go s.readLoop()
}
func (s *Server) readLoop() {
	for {
		var matches []*types.Match
		if err := s.conn.ReadJSON(&matches); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
			continue
		}
		if err := s.producer.ProduceData(context.Background(), matches); err != nil {
			log.Fatal(err)
			continue
		}
	}
}
