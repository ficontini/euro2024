package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/ficontini/euro2024/gateway/store"
	"github.com/ficontini/euro2024/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type UserHandler struct {
	store *store.Store
}

func NewUserHandler(store *store.Store) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.UserParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.Status(http.StatusBadRequest).JSON(errors)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return ErrBadRequest()
	}
	if err := h.store.User.Insert(c.Context(), user); err != nil {
		if errors.Is(err, store.ErrEmailAlreadyInUse) {
			return ErrBadRequestCustomMessage(err.Error())
		}
		return err
	}
	return c.JSON(map[string]string{"user": "registered"})
}

func (h *UserHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var req types.AuthRequest
	if err := c.BodyParser(&req); err != nil {
		return ErrBadRequest()
	}
	if errors := req.Validate(); len(errors) > 0 {
		return c.Status(http.StatusBadRequest).JSON(errors)
	}
	user, err := h.store.User.GetByEmail(c.Context(), req.Email)
	if err != nil {
		return ErrInvalidCredentials()
	}
	if !user.IsPasswordValid(req.Password) {
		return ErrInvalidCredentials()
	}
	auth := types.NewAuth(user.UserID)
	if err := h.store.Auth.Insert(c.Context(), auth); err != nil {
		return err
	}
	token, err := createTokenFromAuth(auth)
	if err != nil {
		return err
	}
	resp := types.AuthResponse{Token: token}
	return c.JSON(resp)
}
func (h *UserHandler) HandlePostSignOut(c *fiber.Ctx) error {
	var (
		auth   = c.Context().UserValue("auth").(*types.Auth)
		filter = &types.AuthFilter{
			UserID:   auth.UserID,
			AuthUUID: auth.AuthUUID,
		}
	)
	if err := h.store.Auth.Delete(c.Context(), filter); err != nil {
		return ErrBadRequest()
	}
	c.Context().SetUserValue("auth", nil)
	return c.JSON(map[string]string{"user": "sign out"})
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
