package repository

import (
	"context"
	"errors"

	"github.com/mendelgusmao/eulabs-api/domain"
	"github.com/mendelgusmao/eulabs-api/domain/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) FetchOne(ctx context.Context, username string) (*model.User, error) {
	tx := r.db.WithContext(ctx)
	user := &model.User{}
	result := tx.First(user, "username = ?", username)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return user, domain.ErrNotFound
		}

		return nil, result.Error
	}

	return user, nil
}

func (r *UserRepository) Create(ctx context.Context, user model.User) error {
	tx := r.db.WithContext(ctx)
	result := tx.Create(&user)

	return result.Error
}
