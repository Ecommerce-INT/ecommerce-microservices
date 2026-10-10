package model

type FavouriteID struct {
	UserID    int    `json:"userId"`
	ProductID int    `json:"productId"`
	LikeDate  string `json:"likeDate"`
}

type FavouriteDto struct {
	UserID    int         `json:"userId"`
	ProductID int         `json:"productId"`
	LikeDate  string      `json:"likeDate"`
	User      interface{} `json:"user,omitempty"`
	Product   interface{} `json:"product,omitempty"`
}

type FavouriteCollectionResponse struct {
	Collection []FavouriteDto `json:"collection"`
}
