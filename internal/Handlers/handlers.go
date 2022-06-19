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

// Get all carts from the shop
func GetAllCartsHandler(c *gin.Context) {
	currency := c.Param("currency")

	carts, err := service.GetAllCarts(currency)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, carts)
}

//Get a cart by id from the shop
func GetCartHandler(c *gin.Context) {
	currency := c.Param("currency")
	cartId := c.Param("cartId")

	cart, err := service.GetCartById(cartId, currency)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())

		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, cart)
}

//Submit a new cart
func AddCartHandler(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var cart structs.Cart
	err := decoder.Decode(&cart)
	if err != nil {
		c.String(http.StatusBadRequest, "request body has wrong format: %s\n", err)

		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	cartID, err := service.AddCart(&cart)
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

	c.String(http.StatusOK, "Successful purchase: %s", cartID)
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

//Update an cart
func UpdateCartHandler(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var cart structs.Cart
	err := decoder.Decode(&cart)
	if err != nil {
		c.String(http.StatusBadRequest, "request body has wrong format: %s\n", err)

		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	cart.Code = c.Param("cartId")

	if err = service.UpdateCart(&cart); err != nil {
		c.String(http.StatusInternalServerError, err.Error())

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusOK, "Cart %s successfully updated", cart.Code)
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

	product.Code = c.Param("productId")

	if err = service.UpdateProduct(&product); err != nil {
		c.String(http.StatusInternalServerError, err.Error())

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusInternalServerError, "Product %s updated succesfully!", product.Code)
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

//Delete an cart
func DeleteCartHandler(c *gin.Context) {
	cartId := c.Param("cartId")

	if err := service.DeleteCart(cartId); err != nil {
		c.String(http.StatusInternalServerError, err.Error())

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusInternalServerError, "cart %s deleted", cartId)
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
