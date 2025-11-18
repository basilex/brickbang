package middleware

import (
	"slices"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"brickbang/internal/service"
	"brickbang/internal/utility"
)

var publicRoutes = map[string][]string{
	"/api/v1/aux":           {fiber.MethodGet},
	"/api/v1/auth/login":    {fiber.MethodPost},
	"/api/v1/auth/register": {fiber.MethodPost},
	"/api/v1/auth/refresh":  {fiber.MethodPost}, // allow refresh as public endpoint
}

// NewAuthMiddleware returns fiber.Handler which uses blacklist service
func NewAuthMiddleware(blacklist service.IBlacklistService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()
		method := c.Method()

		if methods, exists := publicRoutes[path]; exists {
			if slices.Contains(methods, method) {
				return c.Next()
			}
		}

		auth := c.Get("Authorization")
		if auth == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
		}

		parts := strings.Fields(auth)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return fiber.NewError(fiber.StatusUnauthorized, "bad authorization header")
		}

		tokenStr := parts[1]
		claims, err := utility.ParseAccessToken(tokenStr)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token: "+err.Error())
		}

		// check expiry explicitly (jwt lib also does, but double-check)
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now().UTC()) {
			return fiber.NewError(fiber.StatusUnauthorized, "token expired")
		}

		// check blacklist
		isBlocked, err := blacklist.IsBlocked(c.Context(), claims.JTI)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if isBlocked {
			return fiber.NewError(fiber.StatusUnauthorized, "token revoked")
		}

		// attach user and session to locals
		c.Locals("user_id", claims.UserID)
		c.Locals("session_id", claims.SessionID)

		return c.Next()
	}
}
