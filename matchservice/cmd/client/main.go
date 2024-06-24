package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ficontini/euro2024/matchservice/pkg/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var (
		httpAddr = "http://localhost:3003"
		grpcAddr = "localhost:3004"
	)
	svc, err := transport.NewHTTPClient(httpAddr)
	if err != nil {
		log.Fatal(err)
	}
	matches, err := svc.GetLiveMatches(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	for _, m := range matches {
		fmt.Printf("%+v\n", m)
	}

	go func() {
		conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		svc := transport.NewGRPCClient(conn)
		matches, err := svc.GetUpcomingMatches(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		for _, m := range matches {
			fmt.Printf("%+v\n", m)
		}
		matches, err = svc.GetMatchesByTeam(context.Background(), "Deutschland")
		if err != nil {
			log.Fatal(err)
		}
		for _, m := range matches {
			fmt.Printf("%+v\n", m)
		}
	}()
}
