package main

import (
	"assignment2/controllers"
	"assignment2/database"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	db := database.ConnectDB()

	orderController := controllers.NewOrderController(db)

	router.POST("/orders", orderController.CreateOrder)
	router.GET("/orders", orderController.GetOrders)
	router.GET("/orders/:orderid", orderController.GetOrderByID)

	PORT := ":4000"
	router.Run(PORT)
}
