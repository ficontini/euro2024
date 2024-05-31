package matchendpoint

import (
	"context"
	"time"

	"github.com/ficontini/euro2024/matchservice/pkg/service"
	"github.com/ficontini/euro2024/types"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/sony/gobreaker"
)

type Set struct {
	GetUpcomingMatchesEndpoint endpoint.Endpoint
	GetLiveMatchesEndpoint     endpoint.Endpoint
}

func (s Set) GetUpcomingMatches(ctx context.Context) ([]*types.Match, error) {
	resp, err := s.GetUpcomingMatchesEndpoint(ctx, MatchRequest{})
	if err != nil {
		return nil, err
	}
	response := resp.(MatchResponse)
	var matches []*types.Match
	for _, m := range response.Matches {
		matches = append(matches, types.NewMatch(
			m.Date,
			types.NewLocation(m.Location.City, m.Location.Stadium),
			m.Home.Team,
			m.Away.Team,
			m.Status,
			types.NewResult(
				m.Home.Goals,
				m.Away.Goals,
			)))
	}
	return matches, nil
}
func (s Set) GetLiveMatches(ctx context.Context) ([]*types.Match, error) {
	resp, err := s.GetLiveMatchesEndpoint(ctx, struct{}{})
	if err != nil {
		return nil, err
	}
	response := resp.(MatchResponse)
	var matches []*types.Match
	for _, m := range response.Matches {
		matches = append(matches, types.NewMatch(
			m.Date,
			types.NewLocation(m.Location.City, m.Location.Stadium),
			m.Home.Team,
			m.Away.Team,
			m.Status,
			types.NewResult(
				m.Home.Goals,
				m.Away.Goals,
			)))
	}
	return matches, nil
}

func New(svc service.Service) Set {
	var (
		upcomingEndpoint endpoint.Endpoint
		liveEndpoint     endpoint.Endpoint
	)
	{
		upcomingEndpoint = makeGetUpcomingMatchesEndpoint(svc)
		upcomingEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(upcomingEndpoint)
		liveEndpoint = makeGetLiveMatchesEndpoint(svc)
		liveEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(liveEndpoint)
	}
	return Set{
		GetUpcomingMatchesEndpoint: upcomingEndpoint,
		GetLiveMatchesEndpoint:     liveEndpoint,
	}
}
func makeGetLiveMatchesEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		result, err := svc.GetLiveMatches(ctx)
		if err != nil {
			return nil, err
		}
		response := makeResponse(result)
		return response, nil
	}

}
func makeGetUpcomingMatchesEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		result, err := svc.GetUpcomingMatches(ctx)
		if err != nil {
			return nil, err
		}
		response := makeResponse(result)
		return response, nil
	}
}
func makeResponse(matches []*types.Match) MatchResponse {
	var response []Match
	for _, match := range matches {
		response = append(response, Match{
			Date: match.Date,
			Location: Location{
				match.Location.Stadium,
				match.Location.City,
			},
			Status: match.Status,
			Home: Team{
				Team:  match.Home,
				Goals: match.Result.Home,
			},
			Away: Team{
				Team:  match.Away,
				Goals: match.Result.Away,
			},
		})
	}
	return MatchResponse{
		Matches: response,
	}
}

type MatchRequest struct{}
type MatchResponse struct {
	Matches []Match `json:"matches"`
}
type Match struct {
	Home     Team              `json:"home"`
	Away     Team              `json:"away"`
	Date     time.Time         `json:"date"`
	Status   types.MatchStatus `json:"-"`
	Location Location          `json:"location"`
}

type Location struct {
	City    string `json:"city"`
	Stadium string `json:"stadium"`
}

type Team struct {
	Team  string `json:"team"`
	Goals int    `json:"goals"`
}
