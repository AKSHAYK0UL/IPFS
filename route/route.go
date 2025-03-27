package route

import (
	"github.com/gin-gonic/gin"
	"github.com/koulipfs/controller"
)

func RouteTable() *gin.Engine {
	route := gin.Default()
	route.POST("/add-transaction", controller.AddTransactionController)
	route.GET("/transactions", controller.GetTransactionController)
	return route
}
