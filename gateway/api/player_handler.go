package api

import (
	"fmt"

	"github.com/ficontini/euro2024/playerservice/pkg/service"
	"github.com/gofiber/fiber/v2"
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
	param := c.Params("team")
	team, err := validateTeamParameter(param)
	if err != nil {
		return ErrBadRequestCustomMessage("invalid param")
	}
	players, err := h.svc.GetPlayersByTeam(c.Context(), team)
	if err != nil {
		return err
	}
	if len(players) == 0 {
		return ErrResourceNotFound(fmt.Sprintf("players not found for team: %s", team))
	}
	return c.JSON(players)
}
