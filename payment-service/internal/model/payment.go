package model

type OrderDto struct {
	OrderId   *int     `json:"orderId,omitempty"`
	OrderDate string   `json:"orderDate,omitempty"`
	OrderDesc string   `json:"orderDesc,omitempty"`
	OrderFee  *float64 `json:"orderFee,omitempty"`
	ProductId *int     `json:"productId,omitempty"`
}

type UserDto struct {
	ID       *int64 `json:"id,omitempty"`
	FullName string `json:"fullname,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

type PaymentDto struct {
	PaymentId     *int      `json:"paymentId,omitempty"`
	IsPayed       *bool     `json:"isPayed,omitempty"`
	PaymentStatus string    `json:"paymentStatus,omitempty"`
	OrderId       *int      `json:"orderId,omitempty"`
	UserId        *int64    `json:"userId,omitempty"`
	Order         *OrderDto `json:"order,omitempty"`
	User          *UserDto  `json:"user,omitempty"`
}

type KafkaPaymentDto struct {
	PaymentId     *int   `json:"paymentId,omitempty"`
	IsPayed       *bool  `json:"isPayed,omitempty"`
	PaymentStatus string `json:"paymentStatus,omitempty"`
	OrderId       *int   `json:"orderId,omitempty"`
	UserId        *int64 `json:"userId,omitempty"`
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
