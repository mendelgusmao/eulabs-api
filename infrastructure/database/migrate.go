package database

import (
	"log"

	"github.com/mendelgusmao/eulabs-api/domain/model"
)

var models = []any{
	&model.Product{},
	&model.User{},
}

func Migrate(dsn string) error {
	db, err := Connect(dsn)

	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	return db.AutoMigrate(models...)
}
