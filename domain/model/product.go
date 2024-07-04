package model

import (
	"time"
)

type Product struct {
	ID            int64 `gorm:"primarykey"`
	Name          string
	Description   string
	Price         float64
	Quantity      int
	Category      string
	Brand         string
	ImageURL      string
	Weight        float64
	Dimensions    string
	Barcode       string
	SKU           string
	AverageRating float64
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
