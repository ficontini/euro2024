package service

import (
	"context"
	"os"

	"github.com/ficontini/euro2024/playerservice/store"
	"github.com/ficontini/euro2024/types"
	"github.com/go-kit/log"
)

type Service interface {
	GetPlayersByTeam(context.Context, string) ([]*types.Player, error)
}

type basicService struct {
	store store.Storer
}

func newBasicService(store store.Storer) Service {
	return &basicService{
		store: store,
	}
}
func New(store store.Storer) Service {
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
func (svc *basicService) GetPlayersByTeam(ctx context.Context, team string) ([]*types.Player, error) {
	return svc.store.GetPlayersByTeam(ctx, team)
}
