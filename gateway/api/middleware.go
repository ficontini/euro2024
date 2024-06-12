package api

import (
	"github.com/ficontini/euro2024/gateway/service"
	"github.com/gofiber/fiber/v2"
)

func NewAuthMiddleware(svc service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["Authorization"]
		if !ok {
			return ErrUnAuthorized()
		}
		tokenStr := token[0]

		auth, err := svc.Validate(c.Context(), tokenStr)
		if err != nil {
			return ErrUnAuthorized()
		}

		c.Context().SetUserValue("auth", auth)
		return c.Next()
	}
}
