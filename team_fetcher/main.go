package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ficontini/euro2024/team_fetcher/service"
	"github.com/joho/godotenv"
)

const (
	api_key_env  = "API_KEY"
	api_host_env = "API_HOST"
)

func main() {
	var (
		fetcher = service.NewAPIFetcher(
			os.Getenv(api_host_env),
			os.Getenv(api_key_env))
		proccesor = service.NewAPIProcessor()
		svc       = service.New(fetcher, proccesor)
	)
	teams, err := svc.FetchTeams()
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range teams {
		fmt.Printf("%v\n", t.Name)
		for _, p := range t.Players {
			fmt.Printf("%+v\n", p)
		}
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
