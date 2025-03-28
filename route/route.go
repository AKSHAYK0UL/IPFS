package route

import (
	"github.com/gin-gonic/gin"
	"github.com/koulipfs/controller"
	"github.com/koulipfs/middleware"
)

func RouteTable() *gin.Engine {
	route := gin.Default()
	route.POST("/add-transaction", middleware.AuthHeaderMiddleware(), controller.AddTransactionController)
	route.GET("/transactions", controller.GetTransactionController)
	route.GET("/transaction/:id", controller.GetByIdTransactionController)
	return route
}
