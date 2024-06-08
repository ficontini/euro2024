package main

import (
	"context"

	"github.com/ficontini/euro2024/player_storer/store"
	"github.com/ficontini/euro2024/types"
)

type Service interface {
	InsertPlayer(context.Context, *types.Player) error
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
		svc Service
	)
	svc = newBasicService(store)
	svc = newLogMiddleware(svc)
	return svc
}
func (svc *basicService) InsertPlayer(ctx context.Context, player *types.Player) error {
	return svc.store.InsertPlayer(ctx, player)
}
