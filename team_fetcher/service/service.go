package service

import (
	"log"

	"github.com/ficontini/euro2024/types"
)

type Fetcher interface {
	fetchTeams() (*ApiTeamResp, error)
	fetchPlayersByTeam(int) (*ApiPlayerResp, error)
}
type Processor interface {
	processTeam(*TeamResp) (*types.Team, error)
	processPlayers(*ApiPlayerResp) ([]*types.Player, error)
}

type Service interface {
	FetchTeams() ([]*types.Team, error)
}

type basicService struct {
	fetcher   Fetcher
	processor Processor
}

func New(fetcher Fetcher, processor Processor) Service {
	return &basicService{
		fetcher:   fetcher,
		processor: processor,
	}
}

func (svc *basicService) FetchTeams() ([]*types.Team, error) {
	var (
		teams []*types.Team
	)
	res, err := svc.fetcher.fetchTeams()
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range res.Response {
		team, err := svc.processTeam(&r)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	return teams, nil
}
func (svc *basicService) processTeam(r *TeamResp) (*types.Team, error) {
	team, err := svc.processor.processTeam(r)
	if err != nil {
		return nil, err
	}
	players, err := svc.fetchAndProcessPlayers(r.Team.ID)
	if err != nil {
		return nil, err
	}
	team.Players = players
	return team, err
}
func (svc *basicService) fetchAndProcessPlayers(id int) ([]*types.Player, error) {
	res, err := svc.fetcher.fetchPlayersByTeam(id)
	if err != nil {
		return nil, err
	}
	players, err := svc.processor.processPlayers(res)
	if err != nil {
		return nil, err
	}
	return players, nil
}
