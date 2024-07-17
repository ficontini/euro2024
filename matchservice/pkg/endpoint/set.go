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
	GetMatchesByTeamEndpoint   endpoint.Endpoint
	GetEuroWinnerEndpoint      endpoint.Endpoint
}

func (s Set) GetUpcomingMatches(ctx context.Context) ([]*types.Match, error) {
	resp, err := s.GetUpcomingMatchesEndpoint(ctx, MatchRequest{})
	if err != nil {
		return nil, err
	}
	response := resp.(MatchResponse)
	matches := NewMatchesFromMatchResponse(response)
	return matches, nil
}
func (s Set) GetLiveMatches(ctx context.Context) ([]*types.Match, error) {
	resp, err := s.GetLiveMatchesEndpoint(ctx, MatchRequest{})
	if err != nil {
		return nil, err
	}
	response := resp.(MatchResponse)
	matches := NewMatchesFromMatchResponse(response)
	return matches, nil
}
func (s Set) GetMatchesByTeam(ctx context.Context, team string) ([]*types.Match, error) {
	resp, err := s.GetMatchesByTeamEndpoint(ctx, &TeamRequest{Team: team})
	if err != nil {
		return nil, err
	}
	response := resp.(MatchResponse)
	matches := NewMatchesFromMatchResponse(response)
	return matches, nil
}
func (s Set) GetEuroWinner(ctx context.Context) (*types.Match, error) {
	resp, err := s.GetEuroWinnerEndpoint(ctx, MatchRequest{})
	if err != nil {
		return nil, err
	}
	response := resp.(WinnerResponse)
	match := NewMatchFromWinnerResponse(response)
	return match, nil
}

func New(svc service.Service) Set {
	var (
		upcomingEndpoint endpoint.Endpoint
		liveEndpoint     endpoint.Endpoint
		teamEndpoint     endpoint.Endpoint
		winnerEndpoint   endpoint.Endpoint
	)
	{
		upcomingEndpoint = makeGetUpcomingMatchesEndpoint(svc)
		upcomingEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(upcomingEndpoint)
		liveEndpoint = makeGetLiveMatchesEndpoint(svc)
		liveEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(liveEndpoint)
		teamEndpoint = makeGetMatchesByTeamEndpoint(svc)
		teamEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(teamEndpoint)
		winnerEndpoint = makeGetEuroWinnerEndpoint(svc)
		winnerEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(winnerEndpoint)
	}
	return Set{
		GetUpcomingMatchesEndpoint: upcomingEndpoint,
		GetLiveMatchesEndpoint:     liveEndpoint,
		GetMatchesByTeamEndpoint:   teamEndpoint,
		GetEuroWinnerEndpoint:      winnerEndpoint,
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
func makeGetMatchesByTeamEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*TeamRequest)
		result, err := svc.GetMatchesByTeam(ctx, req.Team)
		if err != nil {
			return nil, err
		}
		response := makeResponse(result)
		return response, nil
	}
}
func makeGetEuroWinnerEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		match, err := svc.GetEuroWinner(ctx)
		if err != nil {
			return nil, err
		}
		response := makeWinnerResponse(match)
		return response, nil
	}
}

type MatchRequest struct{}

type TeamRequest struct {
	Team string `json:"team"`
}
type MatchResponse struct {
	Matches []Match `json:"matches"`
}
type Match struct {
	Home     Team              `json:"home"`
	Away     Team              `json:"away"`
	Date     time.Time         `json:"date"`
	Status   types.MatchStatus `json:"-"`
	Location Location          `json:"location"`
	Round    types.RoundStatus `json:"round"`
}

type Location struct {
	City    string `json:"city"`
	Stadium string `json:"stadium"`
}

type Team struct {
	Team  string `json:"team"`
	Goals int    `json:"goals"`
}
type WinnerResponse struct {
	Team  string `json:"team"`
	Final Match  `json:"finalMatch"`
}
