package apifotball

import (
	"github.com/ficontini/euro2024/match_fetcher/processor"
	"github.com/ficontini/euro2024/types"
)

type APIProcessor struct{}

func NewApiProcessor() processor.Processor {
	return &APIProcessor{}
}
func (p *APIProcessor) ProcessData(data any) ([]*types.Match, error) {
	resp := data.(*APIResponse)
	var (
		matchResp []Match
		matches   []*types.Match
	)
	matchResp = resp.Response
	for _, m := range matchResp {
		match := types.NewMatch(
			m.Fixture.Date,
			types.NewLocation(
				m.Fixture.Venue.City,
				m.Fixture.Venue.Name,
			),
			m.Teams.Home.Name,
			m.Teams.Away.Name,
			processStatus(m.Fixture.Status),
			types.NewResult(
				m.Goals.Home,
				m.Goals.Away,
			),
		)
		matches = append(matches, match)
	}
	return matches, nil
}

// TODO:
func processStatus(status Status) types.MatchStatus {
	switch status.Short {
	case "NS":
		return types.NS
	case "LIVE", "HT", "FT", "ET", "PEN_LIVE":
		return types.LIVE
	case "AET", "BREAK", "FT_PEN":
		return types.FINISH
	default:
		return types.UNKNOWN
	}
}
