package main

import (
	"github.com/gin-gonic/gin"
	"github.com/karthihakrishnan/checkoutservice/internal/handler"
)

func main() {
	router := gin.Default()
	//router.GET("/products", controllers.FuncGetProducts)
	//router.GET("/product/:code", controllers.FuncGetProduct)
	//router.POST("/products", controllers.FuncAddProduct)
	router.GET("/product", handler.GetAllProductHandler)
	router.GET("/order", handler.GetAllOrdersHandler)
	router.GET("/product/:productId", handler.GetProductHandler)
	router.GET("/order/:orderId", handler.GetOrderHandler)

	router.POST("/order", handler.AddOrderHandler)     //working
	router.POST("/product", handler.AddProductHandler) //working

	router.PUT("/order/:orderId", handler.UpdateOrderHandler)
	router.PUT("/product/:productId", handler.UpdateProductHandler)

	router.DELETE("/delete/product/:productId", handler.DeleteProductHandler)
	router.DELETE("/delete/order/:orderId", handler.DeleteOrderHandler)
	router.Run("localhost:8083")
}
