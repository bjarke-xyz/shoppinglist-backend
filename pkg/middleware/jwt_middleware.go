package middleware

import (
	"ShoppingList-Backend/app/models"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		jwksUrl := os.Getenv("JWT_JWKS_URL")
		// keycloakUrl := os.Getenv("JWT_KEYCLOAK_URL")

		refreshInterval := time.Hour

		options := keyfunc.Options{
			RefreshInterval: &refreshInterval,
			RefreshErrorHandler: func(err error) {
				log.Printf("There was an error with the jwt.KeyFunc. Error: %v", err)
			},
		}

		jwks, err := keyfunc.Get(jwksUrl, options)
		if err != nil {
			log.Printf("Failed to create JWKS from resource at the given URL. Error: %v", err)
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

		identityUser := models.IdentityUser{ID: claims.Subject}

		c.Locals("user", identityUser)

		return c.Next()
	}
}

// func JWTProtected() func(*fiber.Ctx) error {
// 	config := jwtMiddleware.Config{
// 		SigningKey:   []byte(os.Getenv("JWT_SECRET_KEY")),
// 		ContextKey:   "jwt", // used in private routes
// 		ErrorHandler: jwtError,
// 	}

// 	return jwtMiddleware.New(config)
// }

// func jwtError(c *fiber.Ctx, err error) error {
// 	if err.Error() == "Missing or malformed JWT" {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": true,
// 			"msg":   err.Error(),
// 		})
// 	}

// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 		"error": true,
// 		"msg":   err.Error(),
// 	})
// }
