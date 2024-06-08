package playerendpoint

import (
	"context"

	"github.com/ficontini/euro2024/playerservice/service"
	"github.com/ficontini/euro2024/types"
	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	GetPlayersByTeamEndpoint endpoint.Endpoint
}

func New(svc service.Service) Set {
	return Set{
		GetPlayersByTeamEndpoint: makeGetPlayersByTeamEndpoint(svc),
	}
}

func (s Set) GetPlayersByTeam(ctx context.Context, team string) ([]*types.Player, error) {
	resp, err := s.GetPlayersByTeamEndpoint(ctx, Request{Team: team})
	if err != nil {
		return nil, err
	}
	response := resp.(Response)

	var players []*types.Player
	for _, p := range response.Players {
		performace := types.NewPerformance(p.Shots, p.Goals, p.Assists)
		cards := types.NewCards(p.YellowCards, p.RedCards)
		player := types.NewPlayer(
			p.FirstName,
			p.LastName,
			team,
			p.Age,
			types.NewStatistics(
				performace,
				cards))
		players = append(players, player)
	}

	return players, nil
}
func makeGetPlayersByTeamEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(Request)
		players, err := svc.GetPlayersByTeam(ctx, req.Team)
		if err != nil {
			return nil, err
		}

		var playerResp []Player
		for _, player := range players {
			playerResp = append(playerResp,
				makeResponse(player))
		}
		return Response{
			Players: playerResp,
		}, nil
	}
}

type Request struct {
	Team string
}

type Response struct {
	Players []Player
}
type Player struct {
	FirstName   string
	LastName    string
	Age         int
	Shots       int
	Goals       int
	Assists     int
	YellowCards int
	RedCards    int
}

func makeResponse(player *types.Player) Player {
	return Player{
		FirstName:   player.FirstName,
		LastName:    player.LastName,
		Age:         player.Age,
		Shots:       player.Shots,
		Goals:       player.Goals,
		Assists:     player.Assists,
		YellowCards: player.YellowCards,
		RedCards:    player.RedCards,
	}
}
