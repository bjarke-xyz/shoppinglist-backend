package middleware

import (
	"ShoppingList-Backend/internal/pkg/user"
	"ShoppingList-Backend/pkg/config"
	"fmt"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

func GetAppUser(ctx *fiber.Ctx) user.AppUser {
	appUser := ctx.Locals("user").(user.AppUser)
	return appUser
}

func JWTProtected(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := zap.S()

		refreshInterval := time.Hour

		options := keyfunc.Options{
			RefreshInterval: &refreshInterval,
			RefreshErrorHandler: func(err error) {
				logger.Errorf("There was an error with the jwt.KeyFunc. Error: %v", err)
			},
		}

		jwks, err := keyfunc.Get(cfg.JwtJwksUrl, options)
		if err != nil {
			logger.Errorf("Failed to create JWKS from resource at the given URL. Error: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not create JWKS",
			})
		}

		authHeader := c.Get("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		jwtB64 := strings.TrimSpace(strings.Split(authHeader, "Bearer")[1])
		claims := jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(jwtB64, &claims, jwks.KeyFunc)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to parse the JWT. Error: %v", err),
			})
		}

		if !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		appUser := user.AppUser{ID: claims.Subject}

		c.Locals("user", appUser)

		return c.Next()
	}
}
