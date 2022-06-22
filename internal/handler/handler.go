package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/karthihakrishnan/checkoutservice/internal/service"
	"github.com/karthihakrishnan/checkoutservice/internal/structs"
)

// @Summary Get all products from the shop
// @Tags         Products
// @Produce  application/json
// @Success 200 {string} string	"Successful request"
// @Failure 500 {string} string "Internal server error"
// @Router /product [get]
func GetAllProductHandler(c *gin.Context) {
	currency := c.Param("currency")

	products, err := services.GetAllProducts(currency)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, products)
}

// @Summary Get a product by id from the shop
// @Tags         Products
// @Param   productId	path   string     true  "ID of the product"
// @Produce  application/json
// @Success 200 {string} string	"Successful request"
// @Failure 404 {string} string "Product with such Id not found"
// @Router /product/{productId} [get]
func GetProductHandler(c *gin.Context) {
	currency := c.Param("currency")
	productId := c.Param("productId")

	product, err := services.GetProductById(productId, currency)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())

		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, product)
}

// @Summary Get all orders from the shop
// @Tags         Orders
// @Produce  application/json
// @Success 200 {string} string	"Successful request"
// @Failure 500 {string} string "Internal server error"
// @Router /order [get]
func GetAllOrdersHandler(c *gin.Context) {
	currency := c.Param("currency")

	orders, err := services.GetAllOrders(currency)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, orders)
}

// @Summary Get a order by id from the shop
// @Tags         Orders
// @Param   orderId		path   string     true  "ID of the order"
// @Produce  application/json
// @Success 200 {string} string	"Successful request"
// @Failure 404 {string} string "Order with such Id not found"
// @Router /order/{orderId} [get]
func GetOrderHandler(c *gin.Context) {
	currency := c.Param("currency")
	orderId := c.Param("orderId")

	order, err := services.GetOrderById(orderId, currency)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())

		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, order)
}

// @Summary Submit a new order
// @Tags         Orders
// @Accept   application/json
// @Param   order	body   structs.ExampleOrderRequest	true  "New order details"
// @Produce  application/json
// @Success 200 {string} string	"Successful request"
// @Failure 404 {string} string "Request has wrong format or not enought quantity of a product"
// @Failure 500 {string} string "Internal server error"
// @Router /order [post]
func AddOrderHandler(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var order structs.Order
	err := decoder.Decode(&order)
	if err != nil {
		c.String(http.StatusBadRequest, "request body has wrong format: %s\n", err)

		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	//orderId, err := database.AddOrder(&order)
	orderId, err := services.AddOrder(&order)
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

	c.String(http.StatusOK, "Successful purchase: %s", orderId)
}

// @Summary Add a new product
// @Tags         Products
// @Accept   application/json
// @Param   order	body   structs.ExampleProductRequest	true  "New product details"
// @Produce  application/json
// @Success 200 {string} string	"Successful request"
// @Failure 404 {string} string "Request has wrong format"
// @Failure 500 {string} string "Internal server error"
// @Router /product [post]
func AddProductHandler(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var product structs.Product
	err := decoder.Decode(&product)
	if err != nil {
		c.String(http.StatusBadRequest, "request body has wrong format: %s\n", err)

		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	productId, err := services.AddProduct(&product)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())

		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.String(http.StatusOK, "Product successfully added id: %s", productId)
}

// @Summary Update an order
// @Tags         Orders
// @Accept   application/json
// @Param   orderId		path   string     true  "ID of the order"
// @Param   order	body   structs.ExampleOrderRequest	true  "Updated order details"
// @Produce  application/json
// @Success 200 {string} string	"Successful request"
// @Failure 404 {string} string "Request has wrong format"
// @Failure 500 {string} string "Internal server error"
// @Router /order/{orderId} [put]
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

	if err = services.UpdateOrder(&order); err != nil {
		c.String(http.StatusInternalServerError, err.Error())

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusOK, "Order %s successfully updated", order.ID)
}

// @Summary Update a product
// @Tags         Products
// @Accept   application/json
// @Param   productId	path   string     true  "ID of the product"
// @Param   order	body   structs.ExampleProductRequest	true  "Updated product details"
// @Produce  application/json
// @Success 200 {string} string	"Successful request"
// @Failure 404 {string} string "Request has wrong format"
// @Failure 500 {string} string "Internal server error"
// @Router /product/{productId} [put]
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

	if err = services.UpdateProduct(&product); err != nil {
		c.String(http.StatusInternalServerError, err.Error())

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusInternalServerError, "Product %s updated succesfully!", product.ID)
}

// @Summary Delete a product
// @Tags         Products
// @Param   productId	path   string     true  "ID of the product"
// @Produce  application/json
// @Success 200 {string} string	"Successful request"
// @Failure 404 {string} string "Request has wrong format"
// @Failure 500 {string} string "Internal server error"
// @Router /delete/product/{productId} [delete]
func DeleteProductHandler(c *gin.Context) {
	productId := c.Param("productId")

	if err := services.DeleteProduct(productId); err != nil {
		c.String(http.StatusInternalServerError, err.Error())

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusOK, "Product %s deleted", productId)
}

// @Summary Delete an order
// @Tags         Orders
// @Param   orderId		path   string    true  "ID of the order"
// @Produce  application/json
// @Success 200 {string} string	"Successful request"
// @Failure 404 {string} string "Request has wrong format"
// @Failure 500 {string} string "Internal server error"
// @Router /delete/order/{orderId} [delete]
func DeleteOrderHandler(c *gin.Context) {
	orderId := c.Param("orderId")

	if err := services.DeleteOrder(orderId); err != nil {
		c.String(http.StatusInternalServerError, err.Error())

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusInternalServerError, "order %s deleted", orderId)
}
