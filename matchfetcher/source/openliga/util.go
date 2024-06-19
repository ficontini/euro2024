package openliga

import (
	"time"

	"github.com/ficontini/euro2024/types"
)

func calculateStatus(date time.Time, isFinished bool) types.MatchStatus {
	now := time.Now().UTC()
	if date.After(now) {
		return types.NS
	}

	if isFinished {
		return types.FINISH
	}

	return types.LIVE
}

type Result struct {
	Team1 int
	Team2 int
}

func newResult(team1, team2 int) *Result {
	return &Result{
		Team1: team1,
		Team2: team2,
	}
}

func calculateResult(goals []Goal) *Result {
	if len(goals) == 0 {
		return newResult(0, 0)
	}
	lastGoal := goals[len(goals)-1]
	return newResult(lastGoal.Team1, lastGoal.Team2)
}
