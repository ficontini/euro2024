package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ficontini/euro2024/matchservice/store"
	"github.com/ficontini/euro2024/matchservice/store/fixtures"
	"github.com/ficontini/euro2024/types"
	"github.com/stretchr/testify/assert"
)

func TestHTTPGetLiveMatches(t *testing.T) {
	var (
		store     = store.NewInMemoryStore()
		match     = newLiveMatch(store, "Germany", "Scotland")
		noStarted = newNoStartedMatch(store, "France", "Turkey")
	)
	{
		newLiveMatch(store, "Spain", "Croatia")
	}
	srv := setupHTTPServer(t, store)
	defer srv.Close()

	client, err := setupHTTPClient(t, fmt.Sprintf("%s/matches/live", srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	matches, err := client.GetLiveMatches(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 2, len(matches))
	assert.Equal(t, match.Home.Name, matches[0].Home.Name)
	assert.Equal(t, match.Away.Name, matches[0].Away.Name)
	for _, r := range matches {
		assert.NotEqual(t, r.Home.Name, noStarted.Home.Name)
	}
}
func TestGRPCGetLiveMatches(t *testing.T) {
	var (
		store  = store.NewInMemoryStore()
		match  = newLiveMatch(store, "Germany", "Scotland")
		match1 = newNoStartedMatch(store, "Spain", "Croatia")
	)
	ln, server := setupGRPCServer(t, store)
	defer server.Stop()

	conn, client := setupGRPClient(t, ln)
	defer conn.Close()

	res, err := client.GetLiveMatches(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, match.Home.Name, res[0].Home.Name)
	assert.Equal(t, match.Away.Name, res[0].Away.Name)
	for _, r := range res {
		assert.NotEqual(t, r.Home.Name, match1.Home.Name)
	}
}
func TestGetNoLiveMatches(t *testing.T) {
	var store = store.NewInMemoryStore()
	{
		newNoStartedMatch(store, "Germany", "Scotland")
		newNoStartedMatch(store, "Slovakia", "Denmark")
		newNoStartedMatch(store, "Serbia", "England")
	}
	ln, server := setupGRPCServer(t, store)
	defer server.Stop()

	conn, client := setupGRPClient(t, ln)
	defer conn.Close()

	res, err := client.GetLiveMatches(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, len(res))

}
func TestGetUpcomingMatches(t *testing.T) {
	var (
		store    = store.NewInMemoryStore()
		finished = newFinishedMatch(store, "Portugal", "Czech Republic", types.GROUPS)
		live     = newLiveMatch(store, "France", "Turkey")
	)
	{
		newNoStartedMatch(store, "Germany", "Scotland")
		newNoStartedMatch(store, "Slovakia", "Denmark")
		newNoStartedMatch(store, "Serbia", "England")
	}
	ln, server := setupGRPCServer(t, store)
	defer server.Stop()

	conn, client := setupGRPClient(t, ln)
	defer conn.Close()

	res, err := client.GetUpcomingMatches(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 3, len(res))
	for _, r := range res {
		assert.NotEqual(t, live.Home.Name, r.Home.Name)
		assert.NotEqual(t, finished.Home.Name, r.Home.Name)
	}

}
func TestGetNoUpcomingMatches(t *testing.T) {
	var (
		store   = store.NewInMemoryStore()
		munLoc  = types.NewLocation("Allianz Arena", "München")
		franLoc = types.NewLocation("Deutsche Bank Park", "Frankfurt")
		home1   = types.NewMatchTeam("Serbia", 1)
		away1   = types.NewMatchTeam("England", 3)
		home2   = types.NewMatchTeam("Spain", 1)
		away2   = types.NewMatchTeam("Croatia", 1)
	)
	{
		fixtures.AddMatch(store, time.Now(), munLoc, home1, away1, types.LIVE, types.GROUPS)
		fixtures.AddMatch(store, time.Now(), franLoc, home2, away2, types.LIVE, types.GROUPS)
	}
	ln, server := setupGRPCServer(t, store)
	defer server.Stop()

	conn, client := setupGRPClient(t, ln)
	defer conn.Close()

	res, err := client.GetUpcomingMatches(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, len(res))
}
func TestGetMatchesByTeam(t *testing.T) {
	var (
		store = store.NewInMemoryStore()
		team  = "Germany"
	)
	{
		newFinishedMatch(store, team, "Scotland", types.GROUPS)
		newLiveMatch(store, team, "Hungary")
		newNoStartedMatch(store, "Switzerland", team)
		newNoStartedMatch(store, "Spain", "Italy")
	}
	ln, server := setupGRPCServer(t, store)
	defer server.Stop()

	conn, client := setupGRPClient(t, ln)
	defer conn.Close()

	res, err := client.GetMatchesByTeam(context.Background(), team)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 3, len(res))
	for _, r := range res {
		if r.Home.Name != team && r.Away.Name != team {
			t.Fatalf("it should be %s but it's not", team)
		}
	}
}
func TestGetNoMatchesByTeam(t *testing.T) {
	var (
		store = store.NewInMemoryStore()
		team  = "Portugal"
	)
	{
		newFinishedMatch(store, "Germany", "Scotland", types.GROUPS)
		newLiveMatch(store, "Scotland", "Hungary")
		newNoStartedMatch(store, "Switzerland", "Germany")
		newNoStartedMatch(store, "Spain", "Italy")
	}
	ln, server := setupGRPCServer(t, store)
	defer server.Stop()

	conn, client := setupGRPClient(t, ln)
	defer conn.Close()

	res, err := client.GetMatchesByTeam(context.Background(), team)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, len(res))
}

func TestGetEuroWinner(t *testing.T) {
	var (
		store  = store.NewInMemoryStore()
		winner = "germany"
	)
	newFinishedMatch(store, winner, "england", types.FINAL)

	ln, server := setupGRPCServer(t, store)
	defer server.Stop()

	conn, client := setupGRPClient(t, ln)
	defer conn.Close()

	res, err := client.GetEuroWinner(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, res.Team, winner)
}

func TestGetNoEuroWinner(t *testing.T) {
	var (
		store = store.NewInMemoryStore()
	)
	{
		newFinishedMatch(store, "Germany", "Scotland", types.GROUPS)
		newFinishedMatch(store, "Spain", "Italy", types.GROUPS)
		newFinishedMatch(store, "France", "Austria", types.GROUPS)
	}

	ln, server := setupGRPCServer(t, store)
	defer server.Stop()

	conn, client := setupGRPClient(t, ln)
	defer conn.Close()

	res, _ := client.GetEuroWinner(context.Background())
	assert.Nil(t, res)
}
