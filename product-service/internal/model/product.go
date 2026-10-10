package model

type CategoryDto struct {
	CategoryId     int          `json:"categoryId"`
	CategoryTitle  string       `json:"categoryTitle"`
	ImageUrl       string       `json:"imageUrl"`
	ParentCategory *CategoryDto `json:"parentCategory,omitempty"`
}

type ProductDto struct {
	ProductId    int          `json:"productId"`
	ProductTitle string       `json:"productTitle"`
	ImageUrl     string       `json:"imageUrl"`
	Sku          string       `json:"sku"`
	PriceUnit    float64      `json:"priceUnit"`
	Quantity     int          `json:"quantity"`
	Category     *CategoryDto `json:"category,omitempty"`
}
