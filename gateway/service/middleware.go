package service

import (
	"context"

	"github.com/ficontini/euro2024/types"
	"github.com/sirupsen/logrus"
)

type logMiddleware struct {
	next UserService
}

func newLogMiddleware(next UserService) UserService {
	return &logMiddleware{
		next: next,
	}
}

func (mw *logMiddleware) Insert(ctx context.Context, params types.UserParams) (err error) {
	defer func() {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Info("Insert user")
	}()
	err = mw.next.Insert(ctx, params)
	return err
}
func (mw *logMiddleware) Authenticate(ctx context.Context, req types.AuthRequest) (resp types.AuthResponse, err error) {
	defer func() {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Info("Authenticate user")
	}()
	resp, err = mw.next.Authenticate(ctx, req)
	return resp, err
}
func (mw *logMiddleware) Validate(ctx context.Context, token string) (auth *types.Auth, err error) {
	defer func() {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Info("Validate token")
	}()
	auth, err = mw.next.Validate(ctx, token)
	return auth, err
}
func (mw *logMiddleware) SignOut(ctx context.Context, filter *types.AuthFilter) (err error) {
	defer func() {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Info("Sign out")
	}()
	err = mw.next.SignOut(ctx, filter)
	return err
}
