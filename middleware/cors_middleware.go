package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	corsMiddleware := cors.New(cors.Config{
		AllowOrigins:     []string{"https://zahiqibrahi.github.io"},
		AllowMethods:     []string{http.MethodGet, http.MethodOptions},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})

	return func(ctx *gin.Context) {
		corsMiddleware(ctx)
		ctx.Next()
	}
}
