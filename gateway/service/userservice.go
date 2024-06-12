package service

import (
	"context"
	"fmt"
	"os"

	"github.com/ficontini/euro2024/gateway/store"
	"github.com/ficontini/euro2024/types"
	"github.com/golang-jwt/jwt"
)

type UserService interface {
	Insert(context.Context, types.UserParams) error
	Authenticate(context.Context, types.AuthRequest) (types.AuthResponse, error)
	Validate(context.Context, string) (*types.Auth, error)
	SignOut(context.Context, *types.AuthFilter) error
}

type userService struct {
	store *store.Store
}

func newUserService(store *store.Store) UserService {
	return &userService{
		store: store,
	}
}
func NewUserService(store *store.Store) UserService {
	var svc UserService
	{
		svc = newUserService(store)
		svc = newLogMiddleware(svc)
	}
	return svc
}

func (svc *userService) Insert(ctx context.Context, params types.UserParams) error {
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}
	return svc.store.User.Insert(ctx, user)
}

func (svc *userService) Authenticate(ctx context.Context, req types.AuthRequest) (types.AuthResponse, error) {
	var token types.AuthResponse
	user, err := svc.store.User.GetByEmail(ctx, req.Email)
	if err != nil {
		return token, err
	}
	if !user.IsPasswordValid(req.Password) {
		return token, err
	}
	auth := types.NewAuth(user.UserID)
	if err := svc.store.Auth.Insert(ctx, auth); err != nil {
		return token, err
	}
	tokenStr, err := createTokenFromAuth(auth)
	if err != nil {
		return token, err
	}

	token = types.AuthResponse{Token: tokenStr}

	return token, nil
}

func (svc *userService) Validate(ctx context.Context, tokenStr string) (*types.Auth, error) {
	claims, err := validateToken(tokenStr[len("Bearer "):])
	if err != nil {
		return nil, err
	}
	filter := &types.AuthFilter{
		UserID:   claims["id"].(string),
		AuthUUID: claims["auth_uuid"].(string),
	}
	auth, err := svc.store.Auth.Get(ctx, filter)
	if err != nil {
		return nil, err
	}
	return auth, nil
}

func (svc *userService) SignOut(ctx context.Context, filter *types.AuthFilter) error {
	return svc.store.Auth.Delete(ctx, filter)
}

func createTokenFromAuth(auth *types.Auth) (string, error) {
	claims := jwt.MapClaims{
		"id":        auth.UserID,
		"auth_uuid": auth.AuthUUID,
		"exp":       auth.ExpirationTime,
	}
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return tokenStr, fmt.Errorf("failed to generate auth token")
	}
	return tokenStr, nil
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}
	return claims, nil
}
