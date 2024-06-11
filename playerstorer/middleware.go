package main

import (
	"context"
	"time"

	"github.com/ficontini/euro2024/types"
	"github.com/sirupsen/logrus"
)

type logMiddleware struct {
	next Service
}

func newLogMiddleware(next Service) Service {
	return &logMiddleware{
		next: next,
	}
}
func (mw *logMiddleware) InsertPlayer(ctx context.Context, player *types.Player) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":     time.Since(start),
			"fistName": player.FirstName,
			"lastName": player.LastName,
			"team":     player.Team,
			"err":      err,
		}).Info("InsertPlayer")
	}(time.Now())
	err = mw.next.InsertPlayer(ctx, player)
	return err
}
