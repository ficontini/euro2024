package api

import (
	"fmt"
	"net/http"

	"github.com/ficontini/euro2024/gateway/service"
	"github.com/ficontini/euro2024/types"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
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

	if err := h.svc.Insert(c.Context(), params); err != nil {
		return ErrBadRequest()
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

	token, err := h.svc.Authenticate(c.Context(), req)
	if err != nil {
		return ErrInvalidCredentials()
	}
	if len(token.Token) == 0 {
		return ErrInvalidCredentials()
	}
	fmt.Println("token", token)

	return c.JSON(token)
}

func (h *UserHandler) HandlePostSignOut(c *fiber.Ctx) error {
	var (
		auth   = c.Context().UserValue("auth").(*types.Auth)
		filter = &types.AuthFilter{
			UserID:   auth.UserID,
			AuthUUID: auth.AuthUUID,
		}
	)

	if err := h.svc.SignOut(c.Context(), filter); err != nil {
		return ErrBadRequest()
	}

	c.Context().SetUserValue("auth", nil)

	return c.JSON(map[string]string{"user": "sign out"})
}
