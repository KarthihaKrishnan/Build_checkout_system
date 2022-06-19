package Handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/karthihakrishnan/checkoutservice/internal/service"
	"github.com/karthihakrishnan/checkoutservice/internal/structs"
)

//Get all products from the shop
func GetAllProductHandler(c *gin.Context) {
	currency := c.Param("currency")

	products, err := service.GetAllProducts(currency)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, products)
}

//Get a product by id from the shop
func GetProductHandler(c *gin.Context) {
	currency := c.Param("currency")
	productId := c.Param("productId")

	product, err := service.GetProductById(productId, currency)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())

		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, product)
}

// Get all orders from the shop
func GetAllOrdersHandler(c *gin.Context) {
	currency := c.Param("currency")

	orders, err := service.GetAllOrders(currency)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, orders)
}

//Get a order by id from the shop
func GetOrderHandler(c *gin.Context) {
	currency := c.Param("currency")
	orderId := c.Param("orderId")

	order, err := service.GetOrderById(orderId, currency)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())

		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, order)
}

//Submit a new order
func AddOrderHandler(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var order structs.Order
	err := decoder.Decode(&order)
	if err != nil {
		c.String(http.StatusBadRequest, "request body has wrong format: %s\n", err)

		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	orderID, err := service.AddOrder(&order)
	if err != nil {
		if strings.HasPrefix(err.Error(), "not enough quantity") {
			c.String(http.StatusBadRequest, err.Error())

			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		c.String(http.StatusInternalServerError, err.Error())

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusOK, "Successful purchase: %s", orderID)
}

//Add a new product
func AddProductHandler(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var product structs.Product
	err := decoder.Decode(&product)
	if err != nil {
		c.String(http.StatusBadRequest, "request body has wrong format: %s\n", err)

		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	productID, err := service.AddProduct(&product)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())

		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.String(http.StatusOK, "Product successfully added id: %s", productID)
}

//Update an order
func UpdateOrderHandler(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var order structs.Order
	err := decoder.Decode(&order)
	if err != nil {
		c.String(http.StatusBadRequest, "request body has wrong format: %s\n", err)

		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	order.ID = c.Param("orderId")

	if err = service.UpdateOrder(&order); err != nil {
		c.String(http.StatusInternalServerError, err.Error())

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusOK, "Order %s successfully updated", order.ID)
}

//Update a product
func UpdateProductHandler(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var product structs.Product
	err := decoder.Decode(&product)
	if err != nil {
		c.String(http.StatusBadRequest, "request body has wrong format: %s\n", err)

		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	product.ID = c.Param("productId")

	if err = service.UpdateProduct(&product); err != nil {
		c.String(http.StatusInternalServerError, err.Error())

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusInternalServerError, "Product %s updated succesfully!", product.ID)
}

//Delete a product
func DeleteProductHandler(c *gin.Context) {
	productId := c.Param("productId")

	if err := service.DeleteProduct(productId); err != nil {
		c.String(http.StatusInternalServerError, err.Error())

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusOK, "Product %s deleted", productId)
}

//Delete an order
func DeleteOrderHandler(c *gin.Context) {
	orderId := c.Param("orderId")

	if err := service.DeleteOrder(orderId); err != nil {
		c.String(http.StatusInternalServerError, err.Error())

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusInternalServerError, "order %s deleted", orderId)
}

//	"database/sql"
//	"encoding/json"
//	"net/http"

//	"Build_checkout_system/pkg/models"
//	"Build_checkout_system/pkg/utils"

//	"github.com/jmoiron/sqlx"

// FuncMacBookProPromotion function to check availability of MacBook Pro
/*func FuncMacBookProPromotion(db *sqlx.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Declarations
		respType := utils.ResponseType{
			W: w,
			R: r,
		}
		var err error

		type Result struct {
			Status     int    `json:"status"`
			StatusText string `json:"status_text"`
		}

		type Responses struct {
			Result
			Results interface{}
		}

		type ProductResult struct {
			Scanned_Items string `json:"scanned_item"`
			Total float64 `json:"total"`
		}

		w.Header().Set("Content-Type", "application/json")
		//start processing

		product := models.Product{}
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		err = decoder.Decode(&product)
		if err != nil {
			utils.ErrorResponseHandler("Invalid request object supplied", http.StatusBadRequest, err, respType)
			return
		}

		var tx *sql.Tx
		tx, err = db.Begin()
		if err != nil {
			utils.ErrorResponseHandler("Internal server error 1", http.StatusInternalServerError, err, respType)
			return
		}

		if

	})
} */
