package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/karthihakrishnan/checkoutservice/internal/structs"
	_ "github.com/lib/pq"
	uuid "github.com/nu7hatch/gouuid"
)

const (
	host     = "localhost"
	port     = 5432
	username = "postgres"
	password = "123456"
	dbname   = "postgres"
)

func initDb() *sql.DB {
	log.Println("test")
	psqlconn := fmt.Sprintf("host= %s port = %d user= %s password = %s dbname= %s sslmode=disable", host, port, username, password, dbname)
	log.Println(psqlconn)
	db, err := sql.Open("postgres", psqlconn)
	log.Println(db)
	checkErr(err, "sql.Open failed")

	return db
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func GetAllProducts() ([]structs.Product, error) {
	db := initDb()
	var err error
	if err != nil {
		// simply print the error to the console
		fmt.Println("Err", err.Error())
	}

	defer db.Close()

	products := []structs.Product{}

	rows, err := db.Query("SELECT * FROM products")

	if err != nil {
		return nil, fmt.Errorf("error while reading all products from database: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var prod structs.Product
		if err := rows.Scan(&prod.ID, &prod.Code, &prod.ItemName, &prod.Price, &prod.Qty); err != nil {
			return nil, fmt.Errorf("parsing to a product failed with: %v", err)
		}
		products = append(products, prod)
	}
	return products, nil
}

func GetProductById(productId string) (*structs.Product, error) {

	db := initDb()
	var err error
	prod := &structs.Product{}
	if err != nil {
		// simply print the error to the console
		fmt.Println("Err", err.Error())
	}

	defer db.Close()

	row := db.QueryRow("SELECT * FROM products where id=$1", productId)

	if err := row.Scan(&prod.ID, &prod.Code, &prod.ItemName, &prod.Price, &prod.Qty); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no product with id: %s", productId)
		}
		return nil, fmt.Errorf("searching for %s failed with: %s", productId, err)
	}
	return prod, nil
}

func GetAllOrders() ([]structs.Order, error) {
	db := initDb()
	var err error
	if err != nil {
		// simply print the error to the console
		fmt.Println("Err", err.Error())
	}

	defer db.Close()

	rows, _ := db.Query("SELECT * FROM orders")
	defer rows.Close()

	orders := []structs.Order{}

	for rows.Next() {
		var order structs.Order
		if err := rows.Scan(&order.ID, &order.Code, &order.ItemName, &order.Price, &order.Status); err != nil {
			return nil, fmt.Errorf("getting all products failed with: %v", err)
		}
		products, err := GetAllProductsForOrder(order.ID)
		if err != nil {
			return nil, err
		}
		order.Products = products

		orders = append(orders, order)
	}
	return orders, nil
}

func GetOrderById(orderId string) (*structs.Order, error) {

	db := initDb()
	var err error
	order := &structs.Order{}
	if err != nil {
		// simply print the error to the console
		fmt.Println("Err", err.Error())
	}

	defer db.Close()

	row := db.QueryRow("SELECT * FROM orders where id=$1", orderId)

	if err := row.Scan(&order.ID, &order.Code, &order.ItemName, &order.Price, &order.Status); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no order with id: %s", orderId)
		}
		return nil, fmt.Errorf("searching for %s failed with: %s", orderId, err)
	}

	products, err := GetAllProductsForOrder(order.ID)
	if err != nil {
		return nil, err
	}
	order.Products = products

	return order, nil
}

func AddProduct(product *structs.Product) (string, error) {

	db := initDb()
	var err error

	if err != nil {
		fmt.Println("Err", err.Error())
	}

	// defer the close till after this function has finished
	// executing
	defer db.Close()

	id, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("failed to generate uuid error: %s", err)
	}

	_, err = db.Query("INSERT INTO products (id,code,itemname,price,qty) VALUES ($1,$2,$3,$4,$5)", id.String(), product.Code, product.ItemName, product.Price, product.Qty)

	// if there is an error inserting, handle it
	if err != nil {
		return "", fmt.Errorf("failed to add product to the database, error: %s", err)
	}

	//defer insert.Close()
	return id.String(), nil
}

func AddOrder(order *structs.Order) (string, error) {

	db := initDb()
	var err error

	if err != nil {
		fmt.Println("Err", err.Error())
	}

	// defer the close till after this function has finished
	// executing
	defer db.Close()

	id, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("failed to generate uuid error: %s", err)
	}

	_, err = db.Query("INSERT INTO orders (id,code,itemname,price, status) VALUES ($1,$2,$3,$4, $5)", id.String(), order.Code, order.ItemName, order.Price, order.Status)

	// if there is an error inserting, handle it
	if err != nil {
		return "", fmt.Errorf("failed to add order to the database, error: %s", err)
	}

	//	defer insert.Close()
	return id.String(), nil
}

func UpdateProduct(product *structs.Product) error {
	db := initDb()
	var err error

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}

	// defer the close till after this function has finished
	// executing
	defer db.Close()
	result, err := db.Exec("UPDATE products SET code=$1, itemname=$2, price=$3, qty=$4 WHERE id=$5", product.Code, product.ItemName, product.Price, product.Qty, product.ID)

	if err != nil {
		return fmt.Errorf("failed to update product to the database, error: %s", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no product with code: %s", product.Code)
	}
	return nil
}

func UpdateOrder(order *structs.Order) error {
	db := initDb()
	var err error

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}

	// defer the close till after this function has finished
	// executing
	defer db.Close()
	result, err := db.Exec("UPDATE orders SET itemname=$1, price=$2, code=$3 WHERE id=$4", order.ItemName, order.Price, order.Code, order.ID)

	if err != nil {
		return fmt.Errorf("failed to update order to the database, error: %s", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no order with id: %s", order.ID)
	}
	return nil
}

func DeleteOrder(orderId string) error {
	db := initDb()
	var err error

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}

	// defer the close till after this function has finished
	// executing
	defer db.Close()
	result, err := db.Exec("DELETE FROM orders WHERE id=$1", orderId)

	if err != nil {
		return fmt.Errorf("failed to delete order from the database, error: %s", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no order with id: %s", orderId)
	}
	return nil
}

func DeleteAllProductsForAnOrder(orderId string) error {
	db := initDb()
	var err error

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}

	// defer the close till after this function has finished
	// executing
	defer db.Close()
	result, err := db.Exec("DELETE FROM orderedProduct WHERE order_id = $1;", orderId)
	if err != nil {
		return fmt.Errorf("failed to delete ordered product from the database, error: %s", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no order with id: %s", orderId)
	}

	return nil
}

func DeleteProduct(productId string) error {
	db := initDb()
	var err error

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}

	// defer the close till after this function has finished
	// executing
	defer db.Close()
	result, err := db.Exec("DELETE FROM products WHERE id=$1", productId)

	if err != nil {
		return fmt.Errorf("failed to delete product from the database, error: %s", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no product with id: %s", productId)
	}
	return nil
}

func ChangeProductQuantity(productId string, quantity int) error {
	db := initDb()
	var err error

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}

	// defer the close till after this function has finished
	// executing
	defer db.Close()
	var prod structs.Product

	row := db.QueryRow("SELECT * FROM products WHERE id = $1", productId)
	if err := row.Scan(&prod.Code, &prod.ItemName, &prod.Price, &prod.Qty); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no product with id: %s", productId)
		}
		return fmt.Errorf("searching for id: %s failed with: %s", productId, err)
	}

	newQuantity := prod.Qty - quantity
	if newQuantity < 0 {
		return fmt.Errorf("not enough quantity of product: %s", prod.ItemName)
	}

	if _, err := db.Query("UPDATE products SET qty = $1 WHERE id = $2", newQuantity, prod.ID); err != nil {
		return fmt.Errorf("updating quantity failed with: %s", err)
	}

	return nil
}

func GetAllProductsForOrder(orderId string) ([]structs.Product, error) {
	db := initDb()
	var err error

	if err != nil {
		fmt.Println("Err", err.Error())
	}

	// defer the close till after this function has finished
	// executing
	defer db.Close()
	var products []structs.Product

	rows, err := db.Query("SELECT product_id, quantity FROM orderedProduct WHERE order_id = $1", orderId)
	if err != nil {
		return nil, fmt.Errorf("error while reading ordered product from database: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var prod structs.Product
		if err := rows.Scan(&prod.Code, &prod.Qty); err != nil {
			return nil, fmt.Errorf("parsing to a product failed with: %v", err)
		}

		details, err := GetProductById(prod.ID)
		if err != nil {
			return nil, err
		}
		prod.Code, prod.ItemName, prod.Price, prod.Qty = details.Code, details.ItemName, details.Price, details.Qty

		products = append(products, prod)
	}

	return products, nil
}

func AddOrderedProduct(op *structs.OrderedProduct) error {
	db := initDb()
	var err error

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}

	// defer the close till after this function has finished
	// executing
	defer db.Close()

	id, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("failed to generate uuid error: %s", err)
	}

	_, err = db.Query("INSERT INTO orderedProduct (id, code, product_id, quantity, order_id) VALUES ($1,$2,$3,$4,$5)", id.String(), op.Code, op.ProductId, op.ProductQuantity, op.OrderId)
	if err != nil {
		return fmt.Errorf("failed to add ordered product to the database, error: %s", err)
	}

	//	defer insert.Close()
	return nil
}
