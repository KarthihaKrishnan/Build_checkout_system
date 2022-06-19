package service

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"github.com/karthihakrishnan/checkoutservice/internal/models"
	. "github.com/karthihakrishnan/checkoutservice/internal/structs"
)

type ExchangeRateAPIResponse struct {
	Success   bool   `json:"success"`
	Timestamp int    `json:"timestamp"`
	Base      string `json:"base"`
	Date      string `json:"date"`
	Rates     struct {
		BGN float64 `json:"BGN"`
		CAD float64 `json:"CAD"`
		CHF float64 `json:"CHF"`
		EUR float64 `json:"EUR"`
		GBP float64 `json:"GBP"`
		USD float64 `json:"USD"`
	} `json:"rates"`
}

func GetAllProducts(currency string) ([]Product, error) {
	products, err := models.GetAllProducts()
	if err != nil {
		return nil, fmt.Errorf("failed to get all products with error: %s\n", err)
	}

	for i := range products {
		err := convertPrice(&products[i], currency)
		if err != nil {
			return nil, err
		}
	}

	return products, nil
}

func GetProductById(id string, currency string) (*Product, error) {
	product, err := models.GetProductById(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find such product error: %s\n", err)
	}

	if err = convertPrice(&product, currency); err != nil {
		return nil, err
	}

	return product, nil
}

func GetAllCarts(currency string) ([]Cart, error) {
	carts, err := models.GetAllCarts()
	if err != nil {
		return nil, fmt.Errorf("failed to get all products with error: %s\n", err)
	}

	for i := range carts {
		if err = convertPrice(&carts[i], currency); err != nil {
			return nil, err
		}
	}

	return carts, nil
}

func GetCartById(id string, currency string) (*Cart, error) {
	cart, err := models.GetCartById(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find such cart error: %s\n", err)
	}

	if err = convertPrice(&cart, currency); err != nil {
		return nil, err
	}

	return cart, nil
}

func AddCart(cart *Cart) (string, error) {
	totalPrice := 0.0

	for _, p := range cart.Products {
		err := models.ChangeProductQuantity(p.Code, int(p.Quantity))
		if err != nil {
			return "", err
		}

		product, err := models.GetProductById(p.Code)
		if err != nil {
			return "", err
		}

		totalPrice += product.Price * float64(p.Quantity)
	}

	cart.Price = totalPrice
	cart.Status = "Accepted"

	cartId, err := models.AddCart(cart)
	if err != nil {
		return "", err
	}

	for _, p := range cart.Products {
		err = models.AddCartedProduct(&CartedProduct{
			ProductId:       p.Code,
			ProductQuantity: int(p.Quantity),
			CartId:          cartId,
		})
		if err != nil {
			return "", err
		}
	}

	return cartId, nil
}

func AddProduct(product *Product) (string, error) {
	productId, err := models.AddProduct(product)
	if err != nil {
		return "", err
	}

	return productId, nil
}

func UpdateProduct(product *Product) error {
	err := models.UpdateProduct(product)
	if err != nil {
		return err
	}

	return nil
}

func UpdateCart(cart *Cart) error {
	err := models.UpdateCart(cart)
	if err != nil {
		return err
	}

	return nil
}

func DeleteCart(cartId string) error {
	if err := models.DeleteAllProductsForACart(cartId); err != nil {
		return err
	}

	if err := models.DeleteCart(cartId); err != nil {
		return err
	}

	return nil
}

func DeleteProduct(productId string) error {
	if err := models.DeleteProduct(productId); err != nil {
		return err
	}

	return nil
}

func convertPrice(object interface{}, currency string) error {
	if currency == "" {
		return nil
	}

	rate, err := getRates(currency)
	if err != nil {
		return err
	}

	switch v := object.(type) {
	case *Cart:
		{
			v.Price = math.Round(rate*v.Price*100) / 100
			for i := range v.Products {
				v.Products[i].Price = math.Round(rate*v.Products[i].Price*100) / 100
			}
		}
	case *Product:
		{
			v.Price = math.Round(rate*v.Price*100) / 100
		}
	default:
		return fmt.Errorf("unsupported type")
	}

	return nil
}

func getRates(currency string) (float64, error) {
	const accessKey = "a3d5d57407a65c0b4fa4853c2e5cbe07"
	url := fmt.Sprintf("http://api.exchangeratesapi.io/v1/latest?access_key=%s&format=1", accessKey)

	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("request to exchange rates API failed with error: %s", err)
	}

	decode := json.NewDecoder(resp.Body)
	var exchangeRateResponse ExchangeRateAPIResponse
	err = decode.Decode(&exchangeRateResponse)
	if err != nil {
		return 0, fmt.Errorf("wrong format from exchange rates API, error: %s", err)
	}

	rate := 1.0
	switch currency {
	case "USD":
		rate = exchangeRateResponse.Rates.USD
	case "BGN":
		rate = exchangeRateResponse.Rates.BGN
	case "EUR":
		rate = 1.0
	case "GBP":
		rate = exchangeRateResponse.Rates.GBP
	case "CAD":
		rate = exchangeRateResponse.Rates.CAD
	case "CHF":
		rate = exchangeRateResponse.Rates.CHF
	default:
		return 0, fmt.Errorf("unsupported currency")
	}

	return rate, nil
}
