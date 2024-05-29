package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const queue_url_env_var = "QUEUE_URL"

func main() {
	svc, err := New(os.Getenv(queue_url_env_var))
	if err != nil {
		log.Fatal(err)
	}

	var (
		server = NewSever(svc)
	)

	log.Fatal(server.Start())
}
func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
