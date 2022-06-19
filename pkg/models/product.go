package models

import (
	"Build_checkout_system/pkg/config"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var db *sql.DB

// Define Product structure
type Product struct {
	Code     string  `json:"code"`
	Itemname string  `json:"itemname"`
	Price    float64 `json:"price"`
	Quantity int64   `json:"quantity"`
}

// Define JsonResponse structure
type JsonResponse struct {
	Type    string    `json:"type"`
	Data    []Product `json:"data"`
	Message string    `json:"message"`
}

// Define Cart structure
type Cart struct {
	Code     string  `json:"code"`
	ItemName string  `jspn:"itemname"`
	Price    float64 `json:"price"`
	Quantity int64   `json:"quantity"`
	Status   bool    `json:"status"`
}

// Define Total amount and Shopped items structure
type CartPromo struct {
	Total   float64           `json:"total"`
	Rules   map[string]string `json:"rules"`
	Current []Cart            `json:"current"`
}

// handling messages
func printMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")
}

type allProducts []Product

/*var products = allProducts{
	Product{Item_name: "MacBook Pro", Price: "5399.99", Quantity: "1"},
	Product{Item_name: "Google Home", Price: "49.99", Quantity: "3"},
}*/

// Rules for BuyThreePayTwoOnly
//var Rule_BuyThreePayTwoOnly = map[string]bool{"ult_small": true}

// Set connection info of PostgreSQL Server
func setupDB() *sql.DB {
	//func setupDB() error {
	// connection string
	var DbTransactionUserInfo string = fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", //host, port, user, password, dbname)
		viper.GetString("DbHost"), viper.GetInt("DbPort"), viper.GetString("DbUser"), viper.GetString("DbUserPass"), viper.GetString("DbName"))
	// open database
	var err error
	db, err = sql.Open("postgres", DbTransactionUserInfo)
	config.CheckError(err)

	//Set Variables for sqlx Connection Pool
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)
	db.SetConnMaxLifetime(48 * time.Hour) //2 Day

	//close database
	defer db.Close()

	// check db
	err = db.Ping()
	config.CheckError(err)

	fmt.Println("Connected!")
	return db
	//return nil
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

// create product
/*func createProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct Product
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the product only in order to update")
	}

	json.Unmarshal(reqBody, &newProduct)
	products = append(products, newProduct)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newProduct)
}*/

// Create new product
func FuncCreateNewProduct(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	Item_name := r.FormValue("item_name")
	Price := r.FormValue("price")

	Quantity := r.FormValue("quantity")
	var response = JsonResponse{}

	if Item_name == "" || Price == "" || Quantity == "" {
		response = JsonResponse{Type: "error", Message: "You are missing ID or Item_name or Price or Quantity parameter."}
	} else {
		db := setupDB()
		fmt.Println("Inserting new product into DB")
		fmt.Println("Item_name:" + Item_name)
		fmt.Println("Price:" + Price)
		fmt.Println("Quantity:" + Quantity)
		var lastInsertID int
		/*stmt := `INSERT INTO public.purchased_item(Item_name, Price, Quantity) VALUES ($1, $2, $3) returning id`
		_, err := db.Exec(stmt)*/
		err := db.QueryRow("INSERT INTO public.purchased_item( Item_name, Price, Quantity) VALUES($1, $2, $3) returrning id;",
			Item_name, Price, Quantity).Scan(&lastInsertID)

		config.CheckError(err)
		fmt.Println(err)
		response = JsonResponse{Type: "success", Message: "The product has been inerted successfully!"}
	}

	json.NewEncoder(w).Encode(response)

}

//get one product
/*func FuncgetOneProduct(w http.ResponseWriter, r *http.Request) {
	productID := mux.Vars(r)["id"]

	for _, singleProduct := range products {
		if singleProduct.ID == productID {
			json.NewEncoder(w).Encode(singleProduct)
		}
	}
}*/

// Get all product
/*func FuncGetAllProducts(w http.ResponseWriter, r *http.Request) ([]Product, error) {
	db := setupDB()

	printMessage("Getting all products...")
	rows, err := db.Query("SELECT * FROM public.shop_items")
	config.CheckError(err)
	var products []Product

	// For each product
	for rows.Next() {
		var item_name string
		var price float64
		var quantity int64

		err = rows.Scan(&item_name, &price, &quantity)

		config.CheckError(err)

		products = append(products, Product{Itemname: item_name, Price: price, Quantity: quantity})
	}
	var response = JsonResponse{Type: "success", Data: products}
	json.NewEncoder(w).Encode(response)
	return products, nil
}

/*func FuncGetAllProducts(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(products)
}

//update product
func FuncupdateProduct(w http.ResponseWriter, r *http.Request) {
	productID := mux.Vars(r)["id"]
	var updatedProduct Product

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the product only in order to update")
	}
	json.Unmarshal(reqBody, &updatedProduct)

	for i, singleProduct := range products {
		if singleProduct.ID == productID {
			singleProduct.Item_name = updatedProduct.Item_name
			singleProduct.Price = updatedProduct.Price
			singleProduct.Quantity = updatedProduct.Quantity
			products = append(products[:i], singleProduct)
			json.NewEncoder(w).Encode(singleProduct)
		}
	}
}

func FuncDeleteProduct(w http.ResponseWriter, r *http.Request) {
	productID := mux.Vars(r)["id"]

	for i, singleProduct := range products {
		if singleProduct.ID == productID {
			products = append(products[:i], products[i+1:]...)
			fmt.Fprintf(w, "The product with ID %v has been deleted successfully", productID)
		}
	}
}

// DeleteContactById

/*func FuncDeleteContactById(w http.ResponseWriter, r *http.Request) {
	productID := mux.Vars(r)["id"]

	var response = JsonResponse{}

	if productID == "" {
		response = JsonResponse{Type: "error", Message: "You are missing Id parameter."}
	} else {
		db := config.setupDB()
		fmt.Println("Deleting a product from DB")
		_, err := db.Exec("DELETE FROM public.purchased_item WHERE productID = $1", productID)
		config.CheckError(err)
		response = JsonResponse{Type: "success", Message: "The product has been deleted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}*/

// Get AllCart
/*func FuncGetAllCarts(w http.ResponseWriter, r *http.Request) ([]Cart, error) {
	db := setupDB()

	printMessage("Getting all carts...")
	rows, err := db.Query("SELECT * FROM public.shop_items")
	config.CheckError(err)
	var products []Cart

	// For each product
	for rows.Next() {
		var code string
		var item_name string
		var price float64
		var quantity int64
		var status bool

		err = rows.Scan(&code, &item_name, &price, &quantity, &status)

		config.CheckError(err)

		products = append(products, Cart{Code: code, ItemName: item_name, Price: price, Quantity: quantity, Status: status})
	}
	var response = JsonResponse{Type: "success", Data: products}
	json.NewEncoder(w).Encode(response)
	return products, nil
}*/
