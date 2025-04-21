package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	visitors = make(map[string]*rate.Limiter)
	mutex    sync.Mutex
)

func getLimiter(ip string) *rate.Limiter {
	mutex.Lock()
	defer mutex.Unlock()

	if limiter, exists := visitors[ip]; exists {
		return limiter
	}
	newLimiter := rate.NewLimiter(rate.Every(1*time.Second), 3)
	visitors[ip] = newLimiter
	return newLimiter
}

func RateLimiter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		limiter := getLimiter(ip)
		if !limiter.Allow() {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, "too many request")

		}
		ctx.Next()

	}
}
