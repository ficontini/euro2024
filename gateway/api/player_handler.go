package api

import (
	"context"
	"fmt"

	"github.com/ficontini/euro2024/playerservice/pkg/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PlayerHandler struct {
	svc service.Service
}

func NewPlayerHandler(svc service.Service) *PlayerHandler {
	return &PlayerHandler{
		svc: svc,
	}
}
func (h *PlayerHandler) HandleGetPlayersByTeam(c *fiber.Ctx) error {
	requestID := uuid.NewString()
	ctx := context.WithValue(c.Context(), "requestID", requestID)
	param := c.Params("team")
	team, err := validateTeamParameter(param)
	if err != nil {
		return ErrInvalidParam()
	}
	fmt.Println("----!!!!", ctx.Value("requestID"))
	players, err := h.svc.GetPlayersByTeam(ctx, team)
	if err != nil {
		return err
	}
	if len(players) == 0 {
		return ErrResourceNotFound(fmt.Sprintf("players not found for team: %s", team))
	}
	return c.JSON(players)
}
