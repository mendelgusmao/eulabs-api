package database

import (
	"context"
	"log"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/jinzhu/copier"
	"github.com/mendelgusmao/eulabs-api/domain/dto"
	"github.com/mendelgusmao/eulabs-api/domain/model"
	"github.com/mendelgusmao/eulabs-api/repository"
	"github.com/mendelgusmao/eulabs-api/service"
	"gorm.io/gorm"
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

	if err := PopulateUsers(db); err != nil {
		return err
	}

	if err := PopulateProducts(db); err != nil {
		return err
	}

	return nil
}

func PopulateUsers(db *gorm.DB) error {
	result := db.Delete(&model.User{}, "id > 0")

	if result.Error != nil {
		return result.Error
	}

	users := []dto.CreateUser{
		{
			Name:     "Administrator",
			Username: "admin",
			Admin:    true,
			Password: "admin-password",
		},
		{
			Name:     "Regular User",
			Username: "user",
			Admin:    false,
			Password: "user-password",
		},
	}

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	ctx := context.TODO()

	for _, user := range users {
		userService.Create(ctx, user)

		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func PopulateProducts(db *gorm.DB) error {
	result := db.Delete(&model.Product{}, "id > 0")

	if result.Error != nil {
		return result.Error
	}

	for range 100 {
		fakeProduct := FakeProduct{}
		product := model.Product{}

		if err := faker.FakeData(&fakeProduct); err != nil {
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
