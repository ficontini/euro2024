package service

import (
	"context"
	"fmt"
	"os"

	"github.com/ficontini/euro2024/matchservice/store"
	"github.com/ficontini/euro2024/types"
	"github.com/go-kit/log"
)

type Service interface {
	GetUpcomingMatches(context.Context) ([]*types.Match, error)
	GetLiveMatches(context.Context) ([]*types.Match, error)
	GetMatchesByTeam(context.Context, string) ([]*types.Match, error)
	GetEuroWinner(context.Context) (*types.Match, error)
}

type basicService struct {
	store store.Store
}

func newBasicService(store store.Store) Service {
	return &basicService{
		store: store,
	}
}

func New(store store.Store) Service {
	var (
		logger log.Logger
		svc    Service
	)
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger, "service", "match")
	}
	svc = newBasicService(store)
	svc = newLogMiddleware(logger)(svc)
	return svc
}

func (svc *basicService) GetUpcomingMatches(ctx context.Context) ([]*types.Match, error) {
	matches, err := svc.store.Get(ctx)
	if err != nil {
		return nil, err
	}

	var upcoming []*types.Match
	for _, match := range matches {
		if match.IsUpcoming() {
			upcoming = append(upcoming, match)
		}
	}

	return upcoming, nil
}
func (svc *basicService) GetLiveMatches(ctx context.Context) ([]*types.Match, error) {
	matches, err := svc.store.Get(ctx)
	if err != nil {
		return nil, err
	}

	var live []*types.Match
	for _, match := range matches {
		if match.IsLive() {
			live = append(live, match)
		}
	}

	return live, nil
}
func (svc *basicService) GetMatchesByTeam(ctx context.Context, team string) ([]*types.Match, error) {
	return svc.store.GetMatchesByTeam(ctx, team)
}

func (svc *basicService) GetEuroWinner(ctx context.Context) (*types.Match, error) {
	matches, err := svc.store.GetMatchesByRound(ctx, types.FINAL)
	if err != nil {
		return nil, err
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf("euro winner not found")
	}
	return matches[0], nil
}
