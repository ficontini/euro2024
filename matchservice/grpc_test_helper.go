package main

import (
	"context"
	"log"
	"net"
	"testing"

	matchendpoint "github.com/ficontini/euro2024/matchservice/pkg/endpoint"
	"github.com/ficontini/euro2024/matchservice/pkg/service"
	"github.com/ficontini/euro2024/matchservice/pkg/transport"
	"github.com/ficontini/euro2024/matchservice/proto"
	"github.com/ficontini/euro2024/matchservice/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

func setupGRPCServer(t *testing.T, store store.Store) (*bufconn.Listener, *grpc.Server) {
	var (
		svc           = service.New(store)
		endpoint      = matchendpoint.New(svc)
		matchesServer = transport.NewGRPCServer(endpoint)
		ln            = bufconn.Listen(bufSize)
		server        = grpc.NewServer()
	)
	proto.RegisterMatchesServer(server, matchesServer)
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
