package dto

import (
	"time"
)

type BaseProduct struct {
	Name          string    `json:"name" validate:"required"`
	Description   string    `json:"description"`
	Price         float64   `json:"price" validate:"gt=0"`
	Quantity      int       `json:"quantity"`
	Category      string    `json:"category"`
	Brand         string    `json:"brand"`
	ImageURL      string    `json:"imageUrl"`
	Weight        float64   `json:"weight" validate:"gt=0"`
	Dimensions    string    `json:"dimensions" validate:"required"`
	Barcode       string    `json:"barcode"`
	SKU           string    `json:"sku" validate:"required"`
	AverageRating float64   `json:"averageRating"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Product struct {
	ID int64 `json:"id"`
	BaseProduct
}
