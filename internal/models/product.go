package models

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/karthihakrishnan/checkoutservice/internal/config"
	"github.com/karthihakrishnan/checkoutservice/internal/structs"

	_ "github.com/lib/pq"
)

var db *sql.DB

// Define JsonResponse structure
/*type JsonResponse struct {
	Type    string    `json:"type"`
	Data    []Product `json:"data"`
	Message string    `json:"message"`
}


// handling messages
func printMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")
}

type allProducts []Product*/

// Set connection info of PostgreSQL Server
func setupDB() error {

	const (
		DB_USER     = "postgres"
		DB_PASSWORD = "123456"
		DB_NAME     = "postgres"
	)

	//func setupDB() error {
	// connection string
	// open database
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	var err error
	db, err = sql.Open("postgres", dbinfo)
	config.CheckError(err)

	db.Ping()
	config.CheckError(err)

	fmt.Println("Connected!")
	return nil
}

/*var products = allProducts{
	Product{Item_name: "MacBook Pro", Price: "5399.99", Quantity: "1"},
	Product{Item_name: "Google Home", Price: "49.99", Quantity: "3"},
}*/

// Rules for BuyThreePayTwoOnly
//var Rule_BuyThreePayTwoOnly = map[string]bool{"ult_small": true}

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
/*func FuncCreateNewProduct(w http.ResponseWriter, r *http.Request) {
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
		_, err := db.Exec(stmt)
		err := db.QueryRow("INSERT INTO public.purchased_item( Item_name, Price, Quantity) VALUES($1, $2, $3) returning id",
			Item_name, Price, Quantity).Scan(&lastInsertID)

		config.CheckError(err)
		fmt.Println(err)
		response = JsonResponse{Type: "success", Message: "The product has been inerted successfully!"}
	}

	json.NewEncoder(w).Encode(response)

}*/

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
func GetAllProducts() ([]Product, error) {
	var products []Product

	rows, err := db.Query("SELECT * FROM public.shop_items")
	if err != nil {
		return nil, fmt.Errorf("error while reading all products from database: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.Itemname, &p.Price, &p.Quantity); err != nil {
			return nil, fmt.Errorf("parsing to a product failed with: %v", err)
		}
		products = append(products, p)
	}

	return products, nil
}

func GetProductById(productId string) (*Product, error) {
	row := db.QueryRow("SELECT * FROM public.shop_items WHERE id = ?", productId)

	var p Product
	if err := row.Scan(&p.Code, &p.Itemname, &p.Price, &p.Quantity); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no product with id: %s", productId)
		}
		return nil, fmt.Errorf("searching for %s failed with: %s", productId, err)
	}

	return &p, nil
}

func GetAllOrders() ([]Order, error) {
	var orders []Order

	rows, _ := db.Query("SELECT * FROM orders")
	defer rows.Close()

	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.Name, &o.Address, &o.Phone, &o.Price, &o.Status); err != nil {
			return nil, fmt.Errorf("getting all products failed with: %v", err)
		}

		products, err := GetAllProductsForOrder(o.ID)
		if err != nil {
			return nil, err
		}
		o.Products = products

		orders = append(orders, o)
	}

	return orders, nil
}

func GetOrderById(orderId string) (*Order, error) {
	row := db.QueryRow("SELECT * FROM orders WHERE id = ?", orderId)

	var o Order
	if err := row.Scan(&o.ID, &o.Name, &o.Address, &o.Phone, &o.Price, &o.Status); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no order with id: %s", orderId)
		}
		return nil, fmt.Errorf("searching for %s failed with: %s", orderId, err)
	}

	products, err := GetAllProductsForOrder(o.ID)
	if err != nil {
		return nil, err
	}
	o.Products = products

	return &o, nil
}

func AddProduct(product *Product) (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("failed to generate uuid error: %s", err)
	}

	_, err = db.Query("INSERT INTO products (ID, NAME, CATEGORY, QUANTITY, PRICE) VALUES (?,?,?,?,?)", id.String(), product.Name, product.Category, product.Quantity, product.Price)
	if err != nil {
		return "", fmt.Errorf("failed to add product to the database, error: %s", err)
	}

	return id.String(), nil
}

func AddOrder(order *Order) (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("failed to generate uuid error: %s", err)
	}

	_, err = db.Query("INSERT INTO orders (ID, NAME, Address, Phone, Price, Status) VALUES (?,?,?,?,?,?)", id.String(), order.Name, order.Address, order.Phone, order.Price, order.Status)
	if err != nil {
		return "", fmt.Errorf("failed to add order to the database, error: %s", err)
	}

	return id.String(), nil
}

func UpdateProduct(product *Product) error {

	result, err := db.Exec("UPDATE products SET NAME = ?, CATEGORY = ?, QUANTITY = ?, PRICE = ? WHERE ID = ?", product.Name, product.Category, product.Quantity, product.Price, product.ID)
	if err != nil {
		return fmt.Errorf("failed to update product to the database, error: %s", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no product with id: %s", product.ID)
	}

	return nil
}

func UpdateOrder(order *Order) error {

	result, err := db.Exec("UPDATE orders SET NAME = ?, ADDRESS = ?, PHONE = ?, PRICE = ? WHERE ID = ?", order.Name, order.Address, order.Phone, order.Price, order.ID)
	if err != nil {
		return fmt.Errorf("failed to update order to the database, error: %s", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no product with id: %s", order.ID)
	}

	return nil
}

func DeleteOrder(orderId string) error {
	result, err := db.Exec("DELETE FROM orders WHERE ID = ?;", orderId)
	if err != nil {
		return fmt.Errorf("failed to delete order from the database, error: %s", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no order with id: %d", orderId)
	}

	return nil
}

func DeleteAllProductsForAnOrder(orderId string) error {
	result, err := db.Exec("DELETE FROM orderedProduct WHERE ORDER_ID = ?;", orderId)
	if err != nil {
		return fmt.Errorf("failed to delete ordered product from the database, error: %s", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no order with id: %d", orderId)
	}

	return nil
}

func DeleteProduct(productId string) error {
	result, err := db.Exec("DELETE FROM products WHERE ID = ?;", productId)
	if err != nil {
		return fmt.Errorf("failed to delete product from the database, error: %s", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no product with id: %d", productId)
	}

	return nil
}

func ChangeProductQuantity(productId string, quantity int) error {
	var p Product

	row := db.QueryRow("SELECT * FROM products WHERE id = ?", productId)
	if err := row.Scan(&p.ID, &p.Name, &p.Category, &p.Quantity, &p.Price); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no product with id: %s", productId)
		}
		return fmt.Errorf("searching for id: %s failed with: %s", productId, err)
	}

	newQuantity := p.Quantity - quantity
	if newQuantity < 0 {
		return fmt.Errorf("not enough quantity of product: %s", p.Name)
	}

	if _, err := db.Query("UPDATE products SET quantity = ? WHERE id = ?", newQuantity, p.ID); err != nil {
		return fmt.Errorf("updating quantity failed with: %s", err)
	}

	return nil
}

func GetAllProductsForOrder(orderId string) ([]Product, error) {
	var products []Product

	rows, err := db.Query("SELECT product_id, quantity FROM orderedProduct WHERE order_id = ?", orderId)
	if err != nil {
		return nil, fmt.Errorf("error while reading ordered product from database: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Quantity); err != nil {
			return nil, fmt.Errorf("parsing to a product failed with: %v", err)
		}

		details, err := GetProductById(p.ID)
		if err != nil {
			return nil, err
		}
		p.Name, p.Category, p.Price = details.Name, details.Category, details.Price

		products = append(products, p)
	}

	return products, nil
}

func AddOrderedProduct(op *OrderedProduct) error {
	id, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("failed to generate uuid error: %s", err)
	}

	_, err = db.Query("INSERT INTO orderedProduct (ID, PRODUCT_ID, QUANTITY,  ORDER_ID) VALUES (?,?,?,?)", id.String(), op.ProductId, op.ProductQuantity, op.OrderId)
	if err != nil {
		return fmt.Errorf("failed to add ordered product to the database, error: %s", err)
	}

	return nil
}

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
