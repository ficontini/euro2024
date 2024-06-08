package service

import (
	"context"

	"github.com/ficontini/euro2024/types"
	"github.com/go-kit/log"
)

type Middleware func(Service) Service

type LogMiddleware struct {
	next   Service
	logger log.Logger
}

func newLogMiddleware(logger log.Logger) Middleware {
	return func(s Service) Service {
		return &LogMiddleware{
			next:   s,
			logger: logger,
		}
	}
}
func (mw *LogMiddleware) GetPlayersByTeam(ctx context.Context, team string) (players []*types.Player, err error) {
	defer func() {
		var count int
		if players != nil {
			count = len(players)
		}
		mw.logger.Log("method", "GetPlayersByTeam", "team", team, "count:", count, "err", err)
	}()
	players, err = mw.next.GetPlayersByTeam(ctx, team)
	return players, err
}
