package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/karthihakrishnan/checkoutservice/internal/routes"

	"github.com/gorilla/mux"
)

func main() {
	// Init the mux router
	//r := mux.NewRouter().StrictSlash(true)
	r := mux.NewRouter()
	routes.RegisterProductRoutes(r)
	http.Handle("/", r)
	// serve the app
	fmt.Println("Server running at port 9010")
	log.Fatal(http.ListenAndServe(":9010", r))
}
