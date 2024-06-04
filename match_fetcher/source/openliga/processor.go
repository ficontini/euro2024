package openliga

import (
	"time"

	"github.com/ficontini/euro2024/match_fetcher/processor"
	"github.com/ficontini/euro2024/types"
)

const layout = "2006-01-02T15:04:05"

type APIProcessor struct{}

func NewApiProcessor() processor.Processor {
	return &APIProcessor{}
}
func (p *APIProcessor) ProcessData(data any) ([]*types.Match, error) {
	resp := data.([]*Match)
	var matches []*types.Match
	for _, r := range resp {
		date, err := time.Parse(layout, r.Date)
		if err != nil {
			return nil, err
		}
		matches = append(
			matches,
			types.NewMatch(
				date,
				types.NewLocation(
					r.Location.City,
					r.Location.Stadium,
				),
				r.Team1.Name,
				r.Team2.Name,
				types.NS,
				calculateResult(r.Goals),
			))
	}
	return matches, nil
}

func calculateResult(goals []Goal) *types.Result {
	var (
		home int
		away int
	)
	if len(goals) > 0 {
		home = goals[len(goals)-1].Team1
		away = goals[len(goals)-1].Team2
	}
	return types.NewResult(home, away)
}
