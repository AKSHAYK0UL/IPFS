package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	visitors = make(map[string]*rate.Limiter)
	mutex    sync.Mutex
)

func getVistor(ip string) *rate.Limiter {
	mutex.Lock()
	defer mutex.Unlock()
	limiter, ok := visitors[ip]
	if !ok {
		limiter := rate.NewLimiter(1, 3)
		visitors[ip] = limiter
	}
	return limiter

}

func RateLimiter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		limiter := getVistor(ip)
		if !limiter.Allow() {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, "too many request")

		}
		ctx.Next()

	}
}
