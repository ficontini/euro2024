package api

import (
	"errors"

	"github.com/ficontini/euro2024/gateway/store"
	"github.com/ficontini/euro2024/types"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	store store.UserStorer
}

func NewUserHandler(store store.UserStorer) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.UserParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return ErrBadRequest()
	}
	if err := h.store.Insert(c.Context(), user); err != nil {
		if errors.Is(err, store.ErrEmailAlreadyInUse) {
			ErrBadRequestCustomMessage(err.Error())
		}
	}
	return c.JSON(map[string]string{"user": "registered"})
}
