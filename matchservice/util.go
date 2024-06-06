package main

import (
	"time"

	"github.com/ficontini/euro2024/matchservice/store"
	"github.com/ficontini/euro2024/matchservice/store/fixtures"
	"github.com/ficontini/euro2024/types"
)

func newLiveMatch(store store.Store, team1, team2 string) *types.Match {
	var (
		location = types.NewLocation("Allianz Arena", "München")
		home     = types.NewMatchTeam(team1, 0)
		away     = types.NewMatchTeam(team2, 0)
	)
	return fixtures.AddMatch(store, time.Now(), location, home, away, types.LIVE)
}
func newNoStartedMatch(store store.Store, team1, team2 string) *types.Match {
	var (
		location = types.NewLocation("Olympiastadion", "Berlin")
		home     = types.NewMatchTeam(team1, 0)
		away     = types.NewMatchTeam(team2, 0)
	)
	return fixtures.AddMatch(store, time.Now().AddDate(0, 0, 7), location, home, away, types.NS)
}
func newFinishedMatch(store store.Store, team1, team2 string) *types.Match {
	var (
		location = types.NewLocation("Deutsche Bank Park", "Frankfurt")
		home     = types.NewMatchTeam(team1, 0)
		away     = types.NewMatchTeam(team2, 0)
	)
	return fixtures.AddMatch(store, time.Now().AddDate(0, 0, 7), location, home, away, types.FINISH)
}
