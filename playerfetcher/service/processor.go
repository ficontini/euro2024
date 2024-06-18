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
		statistics = r.Statistics[0]
		performace = newPerfomance(statistics)
		player     = r.Player
	)
	return types.NewPlayer(
		player.FirstName,
		player.LastName,
		statistics.Team.Name,
		statistics.Games.Position,
		player.Age,
		performace,
	)
}
func newPerfomance(s Statistics) *types.Statistics {
	cards := types.NewCards(s.Cards.Yellow, s.Cards.Red)
	accuracy := calculateAccuracy(s.Passes)
	performance := types.NewPerformance(s.Goals.Total, s.Goals.Assists, accuracy)
	return types.NewStatistics(performance, cards)
}

func calculateAccuracy(passes Passes) int {
	if passes.Total == 0 {
		return 0
	}
	return (passes.Accuracy / passes.Total) * 100
}
