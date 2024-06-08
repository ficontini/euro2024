package api

import (
	"fmt"

	"github.com/ficontini/euro2024/matchservice/pkg/service"
	"github.com/gofiber/fiber/v2"
)

type MatchHandler struct {
	svc service.Service
}

func NewMatchHandler(svc service.Service) *MatchHandler {
	return &MatchHandler{
		svc: svc,
	}
}

func (h *MatchHandler) HandleGetUpcomingMatches(c *fiber.Ctx) error {
	matches, err := h.svc.GetUpcomingMatches(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(matches)
}
func (h *MatchHandler) HandleGetLiveMatches(c *fiber.Ctx) error {
	matches, err := h.svc.GetLiveMatches(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(matches)
}
func (h *MatchHandler) HandleGetMatchesByTeam(c *fiber.Ctx) error {
	param := c.Params("team")
	team, err := validateTeamParameter(param)
	if err != nil {
		return ErrInvalidParam()
	}
	matches, err := h.svc.GetMatchesByTeam(c.Context(), team)
	if err != nil {
		return err
	}
	if len(matches) == 0 {
		return ErrResourceNotFound(fmt.Sprintf("matches not found for team: %s", team))
	}
	return c.JSON(matches)
}
