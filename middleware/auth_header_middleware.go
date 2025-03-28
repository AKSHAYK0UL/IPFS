package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func AuthHeaderMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		godotenv.Load() // only used in the local environment
		authHeader := ctx.GetHeader("Authorization")
		authHeaderKey := os.Getenv("AUTHHEADERKEY")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "The Authorization header cannot be null")

		} else if authHeader != authHeaderKey {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "Invalid Authorization header")

		}
		ctx.Next()

	}
}
