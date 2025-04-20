package route

import (
	"github.com/gin-gonic/gin"
	"github.com/koulipfs/controller"
	"github.com/koulipfs/middleware"
)

func RouteTable() *gin.Engine {
	router := gin.Default()
	//CORS middleware
	// router.Use(middleware.RateLimiter(), cors.New(cors.Config{
	// 	AllowOrigins:     []string{"https://zahiqibrahim.github.io"},
	// 	AllowMethods:     []string{http.MethodGet, http.MethodOptions},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))
	//global middleware
	router.Use(middleware.CORS(), middleware.RateLimiter())

	router.POST("/add-transaction", middleware.AuthHeaderMiddleware(), controller.AddTransactionController)
	router.GET("/transactions", controller.GetTransactionController)
	router.GET("/transaction/:id", controller.GetByIdTransactionController)
	router.GET("/smart-contract", controller.SmartContractController)
	return router
}
