package service

import (
	"context"

	"github.com/go-kit/log"

	"github.com/ficontini/euro2024/types"
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

func (m *LogMiddleware) GetUpcomingMatches(ctx context.Context) (matches []*types.Match, err error) {
	defer func() {
		var count int
		if matches != nil {
			count = len(matches)
		}
		m.logger.Log("method", "GetUpcomingMatches", "count:", count, "err", err)
	}()
	matches, err = m.next.GetUpcomingMatches(ctx)
	return matches, err
}

func (m *LogMiddleware) GetLiveMatches(ctx context.Context) (matches []*types.Match, err error) {
	defer func() {
		var count int
		if matches != nil {
			count = len(matches)
		}
		m.logger.Log("method", "GetLiveMatches", "count:", count, "err", err)
	}()
	matches, err = m.next.GetLiveMatches(ctx)
	return matches, err
}

// func (m *LogMiddleware) GetMatchesByTeam(ctx context.Context, team string) (matches []*types.Match, err error) {
// 	defer func() {
// 		var count int
// 		if matches != nil {
// 			count = len(matches)
// 		}
// 		m.logger.Log("method", "GetMatchesByTeam", "team", team, "count:", count, "err", err)
// 	}()
// 	matches, err = m.next.GetMatchesByTeam(ctx, team)
// 	return matches, err
// }
