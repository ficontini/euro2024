package main

import (
	"context"
	"log"
	"net"
	"sync"
	"testing"

	playerendpoint "github.com/ficontini/euro2024/playerservice/pkg/endpoint"
	"github.com/ficontini/euro2024/playerservice/pkg/service"
	"github.com/ficontini/euro2024/playerservice/pkg/transport"
	"github.com/ficontini/euro2024/playerservice/proto"
	"github.com/ficontini/euro2024/playerservice/store"
	"github.com/ficontini/euro2024/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

func setupGRPCServer(t *testing.T, store store.Storer) (*bufconn.Listener, *grpc.Server) {
	var (
		svc           = service.New(store)
		ep            = playerendpoint.New(svc)
		playersServer = transport.NewGRPCServer(ep)
		ln            = bufconn.Listen(bufSize)
		server        = grpc.NewServer()
	)
	proto.RegisterPlayersServer(server, playersServer)
	go func() {
		if err := server.Serve(ln); err != nil {
			log.Fatal(err)
		}
	}()
	return ln, server
}
func setupGRPClient(t *testing.T, ln *bufconn.Listener) (*grpc.ClientConn, service.Service) {
	conn, err := grpc.NewClient("passthrough:whatever", grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
		return ln.Dial()
	}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	client := transport.NewGRPCClient(conn)
	return conn, client
}

type MockStore struct {
	mu      sync.RWMutex
	players map[string][]*types.Player
}

func NewMockStore() *MockStore {
	return &MockStore{
		players: make(map[string][]*types.Player),
	}
}

func (m *MockStore) GetPlayersByTeam(_ context.Context, team string) ([]*types.Player, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.players[team], nil
}
func (m *MockStore) AddPlayer(player *types.Player) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.players[player.Team] = append(m.players[player.Team], player)
}
