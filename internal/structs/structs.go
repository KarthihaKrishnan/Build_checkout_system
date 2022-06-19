package structs

// Define Cart structure
type Cart struct {
	Code     string    `json:"code"`
	ItemName string    `jspn:"itemname"`
	Products []Product `json:"products"`
	Price    float64   `json:"price"`
	Quantity int64     `json:"quantity"`
	Status   string    `json:"status"`
}

// Define Product structure
type Product struct {
	Code     string  `json:"code"`
	Itemname string  `json:"itemname"`
	Price    float64 `json:"price"`
	Quantity int64   `json:"quantity"`
}

type CartedProduct struct {
	ID              string
	ProductId       string
	ProductQuantity int
	CartId          string
}

type ExampleOrderRequest struct {
	Itemname string `default:"MacBook Pro"`
	Products []struct {
		Code     string `default:"43N23P"`
		Quantity int    `default:"2"`
	}
}

type ExampleProductRequest struct {
	Itemname string  `default:"Google Home"`
	Quantity int     `default:"1"`
	Price    float64 `default:"49.99"`
}
