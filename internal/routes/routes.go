package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/karthihakrishnan/checkoutservice/internal/models"

	"github.com/gorilla/mux"
)

var RegisterProductRoutes = func(router *mux.Router) {
	// Route handle & endpoints
	//router.HandleFunc("/", homelink)
	router := gin.Default()
	router.HandleFunc("/product", models.FuncCreateNewProduct).Methods("POST")
	//	router.HandleFunc("/products", models.FuncGetAllProducts).Methods("GET")
	//router.HandleFunc("/products/{id}", FuncgetOneProduct).Methods("GET")
	//router.HandleFunc("/products/{id}", FuncupdateProduct).Methods("PATCH")
	//router.HandleFunc("/products/{id}", FuncDeleteProduct).Methods("DELETE")
	//router.Handle("/product/macbook", controllers.FuncMacBookProPromotion(db)).Methods("GET")
	//router.Handle("/product/google", controllers.FuncGoogleHomesPromotion(db)).Methods("GET")
	//router.Handle("/product/alexa", controllers.FuncAlexaSpeakersPromotion(db)).Methods("GET")
}
