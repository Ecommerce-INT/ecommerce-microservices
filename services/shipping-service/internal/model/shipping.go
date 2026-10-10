package model

import "time"

type OrderItemId struct {
	OrderId   int `json:"orderId"`
	ProductId int `json:"productId"`
}

type OrderItemDto struct {
	OrderId         int        `json:"orderId"`
	ProductId       int        `json:"productId"`
	OrderedQuantity int        `json:"orderedQuantity"`
	CreatedAt       *time.Time `json:"createdAt,omitempty"`
	UpdatedAt       *time.Time `json:"updatedAt,omitempty"`
}

type DtoCollectionResponse[T any] struct {
	Collection []T `json:"collection"`
}
