package service

import (
	"context"
	"os"
	"time"

	"github.com/ficontini/euro2024/types"
	"github.com/go-kit/log"
)

type Service interface {
	ProcessData(context.Context, []*types.Match) error
	GetUpcomingMatches(context.Context) ([]*types.Match, error)
	GetLiveMatches(context.Context) ([]*types.Match, error)
	GetMatchesByTeam(context.Context, string) ([]*types.Match, error)
}

type basicService struct {
	store Store
}

func newBasicService() Service {
	return &basicService{
		store: NewInMemoryStore(),
	}
}

func New() Service {
	var (
		logger log.Logger
		svc    Service
	)
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger, "service", "match")
	}
	svc = newBasicService()
	svc = newLogMiddleware(logger)(svc)
	return svc
}

func (svc *basicService) ProcessData(ctx context.Context, matches []*types.Match) error {
	svc.store.Clean(ctx)
	for _, m := range matches {
		svc.store.Add(ctx, m)
	}
	return nil
}

func (svc *basicService) GetUpcomingMatches(ctx context.Context) ([]*types.Match, error) {
	time.Sleep(2 * time.Second)
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
	return nil, nil
}
