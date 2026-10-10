package model

type Inventory struct {
	ID          int64  `json:"id"`
	ProductName string `json:"productName"`
	Quantity    int    `json:"quantity"`
}

type InventoryResponse struct {
	ProductName string `json:"productName"`
	IsInStock   bool   `json:"isInStock"`
	Quantity    int    `json:"quantity,omitempty"`
}
