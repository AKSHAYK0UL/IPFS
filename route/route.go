package route

import (
	"github.com/gin-gonic/gin"
	"github.com/koulipfs/controller"
	"github.com/koulipfs/middleware"
)

func RouteTable() *gin.Engine {
	router := gin.Default()
	// middleware's
	router.Use(middleware.RateLimiter(), middleware.CORS())

	router.POST("/add-transaction", middleware.AuthHeaderMiddleware(), controller.AddTransactionController)
	router.GET("/transactions", controller.GetTransactionController)
	router.GET("/transaction/:id", controller.GetByIdTransactionController)
	return router
}
