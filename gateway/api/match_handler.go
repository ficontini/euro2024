package api

import (
	"errors"
	"strings"

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
		return ErrResourceNotFound(team)
	}
	return c.JSON(matches)
}

func validateTeamParameter(value string) (string, error) {
	if len(value) == 0 {
		return "", errors.New("empty value")
	}
	if len(value) == 1 {
		return "", errors.New("invalid value")
	}
	var (
		builder strings.Builder
		first   = strings.ToUpper(string(value[0]))
		last    = strings.ToLower(string(value[1:]))
	)
	builder.WriteString(first)
	builder.WriteString(last)
	return builder.String(), nil
}
