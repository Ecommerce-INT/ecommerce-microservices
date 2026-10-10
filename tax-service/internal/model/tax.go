package model

import "time"

type TaxClass struct {
	ID             int64      `json:"id"`
	Name           string     `json:"name"`
	CreatedBy      *string    `json:"createdBy,omitempty"`
	CreatedOn      *time.Time `json:"createdOn,omitempty"`
	LastModifiedBy *string    `json:"lastModifiedBy,omitempty"`
	LastModifiedOn *time.Time `json:"lastModifiedOn,omitempty"`
}

type TaxClassVm struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type TaxClassPostVm struct {
	Name string `json:"name"`
}

type TaxClassListGetVm struct {
	TaxClassContent []TaxClassVm `json:"taxClassContent"`
	PageNo          int          `json:"pageNo"`
	PageSize        int          `json:"pageSize"`
	TotalElements   int          `json:"totalElements"`
	TotalPages      int          `json:"totalPages"`
	IsLast          bool         `json:"isLast"`
}

type TaxRate struct {
	ID                int64      `json:"id"`
	Rate              float64    `json:"rate"`
	ZipCode           *string    `json:"zipCode"`
	TaxClassID        int64      `json:"taxClassId"`
	TaxClassName      string     `json:"taxClassName"`
	StateOrProvinceID *int64     `json:"stateOrProvinceId"`
	CountryID         int64      `json:"countryId"`
	CreatedBy         *string    `json:"createdBy,omitempty"`
	CreatedOn         *time.Time `json:"createdOn,omitempty"`
	LastModifiedBy    *string    `json:"lastModifiedBy,omitempty"`
	LastModifiedOn    *time.Time `json:"lastModifiedOn,omitempty"`
}

type TaxRateVm struct {
	ID                int64   `json:"id"`
	Rate              float64 `json:"rate"`
	ZipCode           *string `json:"zipCode"`
	TaxClassId        int64   `json:"taxClassId"`
	TaxClassName      string  `json:"taxClassName"`
	StateOrProvinceId *int64  `json:"stateOrProvinceId"`
	CountryId         int64   `json:"countryId"`
}

type TaxRatePostVm struct {
	Rate              float64 `json:"rate"`
	ZipCode           *string `json:"zipCode"`
	TaxClassId        int64   `json:"taxClassId"`
	StateOrProvinceId *int64  `json:"stateOrProvinceId"`
	CountryId         int64   `json:"countryId"`
}

type TaxRateListGetVm struct {
	TaxRateContent []TaxRateVm `json:"taxRateContent"`
	PageNo         int         `json:"pageNo"`
	PageSize       int         `json:"pageSize"`
	TotalElements  int         `json:"totalElements"`
	TotalPages     int         `json:"totalPages"`
	IsLast         bool        `json:"isLast"`
}
