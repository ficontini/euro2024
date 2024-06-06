package service

import (
	"github.com/ficontini/euro2024/types"
)

type Processor interface {
	ProcessData([]*ApiResponse) []*types.Player
}
type APIProcessor struct {
	players []*types.Player
}

func NewAPIProcessor() Processor {
	return &APIProcessor{}
}

func (p *APIProcessor) ProcessData(resp []*ApiResponse) []*types.Player {
	for _, r := range resp {
		p.processApiResponse(r.Response)
	}
	return p.players
}

func (p *APIProcessor) processApiResponse(resp []PlayerResp) {
	for _, r := range resp {
		p.players = append(p.players, newPlayer(r))
	}
}
func newPlayer(r PlayerResp) *types.Player {
	var (
		performace = newPerfomance(r.Statistics[0])
		player     = r.Player
	)
	return types.NewPlayer(
		player.FirstName,
		player.LastName,
		player.Nationality,
		player.Age,
		performace,
	)
}
func newPerfomance(s Statistics) *types.Statistics {
	var (
		cards       = types.NewCards(s.Cards.Yellow, s.Cards.Red)
		performance = types.NewPerformance(s.Shots.Total, s.Goals.Total, s.Goals.Assists)
	)
	return types.NewStatistics(performance, cards)
}
