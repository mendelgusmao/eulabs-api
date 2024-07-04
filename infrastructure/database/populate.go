package database

import (
	"log"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/jinzhu/copier"
	"github.com/mendelgusmao/eulabs-api/domain/model"
)

type FakeProduct struct {
	Name          string
	Description   string
	Price         float64
	Quantity      int
	Category      string
	Brand         string
	DateAdded     time.Time
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

func Populate(dsn string) error {
	db, err := Connect(dsn)

	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	result := db.Delete(&model.Product{}, "id > 0")

	if result.Error != nil {
		return result.Error
	}

	for range 100 {
		fakeProduct := FakeProduct{}
		product := model.Product{}

		if err := faker.FakeData(&fakeProduct); err != nil {
			log.Println(fakeProduct)
			return err
		}

		if err := copier.Copy(&product, &fakeProduct); err != nil {
			return err
		}

		result := db.Create(&product)

		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}
