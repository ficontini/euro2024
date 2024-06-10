package main

import (
	"context"
	"testing"

	"github.com/ficontini/euro2024/types"
	"github.com/stretchr/testify/assert"
)

func TestGetPlayersByTeam(t *testing.T) {
	var (
		store  = NewMockStore()
		player types.Player
		team   = "France"
	)
	for i := 0; i < 5; i++ {
		player.Team = team
		store.AddPlayer(&player)
	}
	ln, srv := setupGRPCServer(t, store)
	defer srv.Stop()

	conn, client := setupGRPClient(t, ln)
	defer conn.Close()

	res, err := client.GetPlayersByTeam(context.Background(), team)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 5, len(res))
}
func TestGetNoPlayerByTeam(t *testing.T) {
	var (
		store  = NewMockStore()
		player types.Player
		team   = "France"
	)
	for i := 0; i < 5; i++ {
		player.Team = team
		store.AddPlayer(&player)
	}
	ln, srv := setupGRPCServer(t, store)
	defer srv.Stop()

	conn, client := setupGRPClient(t, ln)
	defer conn.Close()

	res, err := client.GetPlayersByTeam(context.Background(), "Portugal")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 0, len(res))
}
