package fixtures

import (
	"context"
	"log"
	"time"

	"github.com/ficontini/euro2024/matchservice/store"
	"github.com/ficontini/euro2024/types"
)

func AddMatch(store store.Store, date time.Time, location *types.Location, home, away *types.MatchTeam, status types.MatchStatus, round types.RoundStatus) *types.Match {
	match := types.NewMatch(date, location, home, away, status, round)
	if err := store.Add(context.Background(), match); err != nil {
		log.Fatal(err)
	}
	return match
}
