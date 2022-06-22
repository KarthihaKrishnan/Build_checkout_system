package structs

type Product struct {
	ID       string
	Code     string
	ItemName string
	Price    float64
	Qty      int
}

type Order struct {
	ID       string    `json:"id"`
	Code     string    `json:"code"`
	ItemName string    `json:"itemname"`
	Products []Product `json:"products"`
	Price    float64   `json:"price"`
	Status   string    `json:"status"`
}

type OrderedProduct struct {
	ID              string
	Code            string
	ProductId       string
	ProductQuantity int
	OrderId         string
}
