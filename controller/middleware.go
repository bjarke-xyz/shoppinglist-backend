package controller

import (
	"ShoppingList-Backend/model"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	bearer = "Bearer"
)

// https://stackoverflow.com/a/65357906
func JWTAuthorized(jwksURL string, keycloakUrl string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		refreshInterval := time.Hour
		options := keyfunc.Options{
			RefreshInterval: &refreshInterval,
			RefreshErrorHandler: func(err error) {
				log.Printf("There was an error with the jwt.KetFunc\nError: %s", err.Error())
			},
		}

		jwks, err := keyfunc.Get(jwksURL, options)
		if err != nil {
			log.Printf("Failed to create JWKS from resource at the given URL.\nError: %s", err.Error())
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, HttpError{
				Code:  http.StatusUnauthorized,
				Error: "Could not create JWKS",
			})
			return
		}

		authHeader := ctx.GetHeader("Authorization")
		if !strings.Contains(authHeader, bearer) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, HttpError{
				Code:  http.StatusUnauthorized,
				Error: "Bearer authentication scheme must be used",
			})
			return
		}

		jwtB64 := strings.TrimSpace(strings.Split(authHeader, "Bearer")[1])
		claims := jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(jwtB64, &claims, jwks.KeyFunc)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, HttpError{
				Code:  http.StatusUnauthorized,
				Error: fmt.Sprintf("Failed to parse the JWT. Error: %s", err.Error()),
			})
			return
		}

		if !token.Valid {
			log.Printf("The token is not valid")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, HttpError{
				Code:  http.StatusUnauthorized,
				Error: "Invalid token",
			})
		}

		identityUser := model.IdentityUser{ID: claims.Subject}

		ctx.Set("user", identityUser)
		ctx.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT, PATCH, DELETE")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next()
	}
}
