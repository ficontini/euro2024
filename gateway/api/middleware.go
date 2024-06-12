package api

import (
	"fmt"
	"os"

	"github.com/ficontini/euro2024/gateway/store"
	"github.com/ficontini/euro2024/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func NewAuthMiddleware(store store.AuthStorer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["Authorization"]
		if !ok {
			return ErrUnAuthorized()
		}
		tokenStr := token[0]
		claims, err := validateToken(tokenStr[len("Bearer "):])
		if err != nil {
			return ErrUnAuthorized()
		}
		filter := &types.AuthFilter{
			UserID:   claims["id"].(string),
			AuthUUID: claims["auth_uuid"].(string),
		}
		auth, err := store.Get(c.Context(), filter)
		if err != nil {
			return ErrUnAuthorized()
		}
		c.Context().SetUserValue("auth", auth)
		return c.Next()
	}
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
