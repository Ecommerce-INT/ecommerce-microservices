package model

type UserDto struct {
	ID       *int64  `json:"id,omitempty"`
	FullName string  `json:"fullname,omitempty"`
	Username string  `json:"username,omitempty"`
	Email    string  `json:"email,omitempty"`
}

type ProductDto struct {
	ProductId   *int     `json:"productId,omitempty"`
	ProductTitle string  `json:"productTitle,omitempty"`
	ImageUrl     string  `json:"imageUrl,omitempty"`
	Sku          string  `json:"sku,omitempty"`
	PriceUnit    *float64 `json:"priceUnit,omitempty"`
	Quantity     *int     `json:"quantity,omitempty"`
}

type OrderDto struct {
	OrderId   *int        `json:"orderId,omitempty"`
	OrderDate string      `json:"orderDate,omitempty"`
	OrderDesc string      `json:"orderDesc,omitempty"`
	OrderFee  *float64    `json:"orderFee,omitempty"`
	ProductId *int        `json:"productId,omitempty"`
	Product   *ProductDto `json:"product,omitempty"`
	Cart      *CartDto    `json:"cart,omitempty"`
}

type CartDto struct {
	CartId *int        `json:"cartId,omitempty"`
	UserId *int64      `json:"userId,omitempty"`
	Orders []OrderDto  `json:"order,omitempty"`
	User   *UserDto    `json:"user,omitempty"`
}

type PageResponse[T any] struct {
	Content          []T   `json:"content"`
	TotalElements    int64 `json:"totalElements"`
	TotalPages       int   `json:"totalPages"`
	Size             int   `json:"size"`
	Number           int   `json:"number"`
	NumberOfElements int   `json:"numberOfElements"`
	First            bool  `json:"first"`
	Last             bool  `json:"last"`
	Empty            bool  `json:"empty"`
}
