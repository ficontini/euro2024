package queue

import (
	"context"

	"github.com/ficontini/euro2024/matchservice/store"
	"github.com/ficontini/euro2024/types"
)

type Service interface {
	ProcessData(context.Context, []*types.Match) error
}
type basicService struct {
	store store.Store
}

func New(store store.Store) Service {
	return &basicService{store: store}
}

func (svc *basicService) ProcessData(ctx context.Context, matches []*types.Match) error {
	svc.store.Clean(ctx)
	for _, m := range matches {
		svc.store.Add(ctx, m)
	}
	return nil
}
