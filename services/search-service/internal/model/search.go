package model

import "time"

type ProductDoc struct {
	ID                    int64     `json:"id"`
	Name                  string    `json:"name"`
	Slug                  string    `json:"slug"`
	Price                 float64   `json:"price"`
	IsPublished           bool      `json:"isPublished"`
	IsVisibleIndividually bool      `json:"isVisibleIndividually"`
	IsAllowedToOrder      bool      `json:"isAllowedToOrder"`
	IsFeatured            bool      `json:"isFeatured"`
	ThumbnailMediaID      *int64    `json:"thumbnailMediaId,omitempty"`
	Brand                 string    `json:"brand"`
	Categories            []string  `json:"categories"`
	Attributes            []string  `json:"attributes"`
	CreatedOn             time.Time `json:"createdOn"`
}

type ProductGetVm struct {
	ID                    int64      `json:"id"`
	Name                  string     `json:"name"`
	Slug                  string     `json:"slug"`
	ThumbnailMediaID      *int64     `json:"thumbnailId,omitempty"`
	Price                 float64    `json:"price"`
	IsAllowedToOrder      bool       `json:"isAllowedToOrder"`
	IsPublished           bool       `json:"isPublished"`
	IsFeatured            bool       `json:"isFeatured"`
	IsVisibleIndividually bool       `json:"isVisibleIndividually"`
	CreatedOn             *time.Time `json:"createdOn,omitempty"`
}

type ProductListGetVm struct {
	Products      []ProductGetVm            `json:"products"`
	PageNo        int                       `json:"pageNo"`
	PageSize      int                       `json:"pageSize"`
	TotalElements int64                     `json:"totalElements"`
	TotalPages    int                       `json:"totalPages"`
	IsLast        bool                      `json:"isLast"`
	Aggregations  map[string]map[string]int `json:"aggregations"`
}

type ProductNameGetVm struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ProductNameListVm struct {
	ProductNames []ProductNameGetVm `json:"productNames"`
}
