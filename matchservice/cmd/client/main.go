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
	// httpAddr := "http://localhost:3003"
	// svc, err := transport.NewHTTPClient(httpAddr)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// matches, err := svc.GetUpcomingMatches(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for _, m := range matches {
	// 	fmt.Printf("%+vn", m)
	// }
	var (
		endpoint = "localhost:3004"
	)
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
		fmt.Printf("%+vn", m)
	}
}
