package openliga

import (
	"time"

	"github.com/ficontini/euro2024/matchfetcher/service"
	util "github.com/ficontini/euro2024/matchfetcher/source"
	"github.com/ficontini/euro2024/types"
)

const (
	layout    = "2006-01-02T15:04:05"
	layoutUTC = "2006-01-02T15:04:05Z"
)

type APIProcessor struct {
	translator Translator
}

func NewApiProcessor() (service.Processor, error) {
	translator, err := NewTranslator()
	if err != nil {
		return nil, err
	}
	return &APIProcessor{
		translator: translator,
	}, nil
}

func (p *APIProcessor) ProcessData(data any) ([]*types.Match, error) {
	resp := data.([]*Match)
	var matches []*types.Match
	for _, r := range resp {
		match, err := p.newMatch(r)
		if err != nil {
			return nil, err
		}
		if len(match.Home.Name) > 0 && len(match.Away.Name) > 0 {
			matches = append(matches, match)
		}
	}
	return matches, nil
}

func (p *APIProcessor) newMatch(m *Match) (*types.Match, error) {
	dateUTC, err := time.Parse(layoutUTC, m.DateUTC)
	if err != nil {
		return nil, err
	}
	date, err := time.Parse(layout, m.Date)
	if err != nil {
		return nil, err
	}

	result := calculateResult(m.Goals)

	location := types.NewLocation(m.Location.City, m.Location.Stadium)
	home := types.NewMatchTeam(p.translator.Get(m.Team1.Name), result.Team1)
	away := types.NewMatchTeam(p.translator.Get(m.Team2.Name), result.Team2)

	status := calculateStatus(dateUTC.UTC(), m.IsFinished)
	return types.NewMatch(
		date,
		location,
		home,
		away,
		status,
		util.GetRound(p.translator.Get(m.Group.Name)),
	), nil

}
