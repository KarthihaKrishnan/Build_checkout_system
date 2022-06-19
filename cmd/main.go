package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/karthihakrishnan/checkoutservice/internal/Handlers"
)

func main() {
	// Init the mux router
	//r := mux.NewRouter().StrictSlash(true)
	/*r := mux.NewRouter()
	routes.RegisterProductRoutes(r)
	http.Handle("/", r)*/
	r := gin.Default()
	r.GET("/product", Handlers.GetAllProductHandler)
	r.GET("/cart", Handlers.GetAllCartsHandler)
	r.GET("/product/:productId", Handlers.GetProductHandler)
	r.GET("/cart/:cartId", Handlers.GetCartHandler)

	r.POST("/cart", Handlers.AddCartHandler)
	r.POST("/product", Handlers.AddProductHandler)

	r.PUT("/cart/:cartId", Handlers.UpdateCartHandler)
	r.PUT("/product/:productId", Handlers.UpdateProductHandler)

	r.DELETE("/delete/product/:productId", Handlers.DeleteProductHandler)
	r.DELETE("/delete/cart/:cartId", Handlers.DeleteCartHandler)

	log.Println("Listening to port 8080...")
	log.Fatal(r.Run(":8080"))
	// serve the app
	//	fmt.Println("Server running at port 9010")
	//	log.Fatal(http.ListenAndServe(":9010", r))
}
