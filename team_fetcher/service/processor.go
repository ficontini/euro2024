package service

import "github.com/ficontini/euro2024/types"

type APIProcessor struct{}

func NewAPIProcessor() Processor {
	return &APIProcessor{}
}

func (p *APIProcessor) processTeam(teamResp *TeamResp) (*types.Team, error) {
	var (
		team    = teamResp.Team
		venue   = teamResp.Venue
		stadium = types.NewStadium(venue.Name, venue.City)
	)
	return types.NewTeam(team.ID, team.Name, stadium), nil
}

func (p *APIProcessor) processPlayers(resp *ApiPlayerResp) ([]*types.Player, error) {
	var players []*types.Player
	for _, p := range resp.Response {
		player := newPlayer(p)
		players = append(players, player)
	}
	return players, nil
}
func newPlayer(p PlayerResp) *types.Player {
	var (
		performace = newPerfomance(p.Statistics[0])
		player     = p.Player
	)
	return types.NewPlayer(
		player.FirstName,
		player.LastName,
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
