package repository

import (
	"context"
	"errors"

	"github.com/mendelgusmao/eulabs-api/domain"
	"github.com/mendelgusmao/eulabs-api/domain/model"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) FetchMany(ctx context.Context, conditions ...any) ([]model.Product, error) {
	tx := r.db.WithContext(ctx)
	products := make([]model.Product, 0)
	result := tx.Find(&products, conditions...)

	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}

func (r *ProductRepository) FetchOne(ctx context.Context, id int64) (*model.Product, error) {
	tx := r.db.WithContext(ctx)
	product := &model.Product{}
	result := tx.First(product, []int64{id})

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return product, domain.ErrNotFound
		}

		return nil, result.Error
	}

	return product, nil
}

func (r *ProductRepository) Create(ctx context.Context, product model.Product) (*model.Product, error) {
	tx := r.db.WithContext(ctx)
	result := tx.Create(&product)

	return &product, result.Error
}

func (r *ProductRepository) Update(ctx context.Context, product model.Product) (*model.Product, error) {
	tx := r.db.WithContext(ctx)
	result := tx.Save(&product)

	if result.RowsAffected == 0 {
		return nil, domain.ErrNotFound
	}

	return &product, nil
}

func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	tx := r.db.WithContext(ctx)
	result := tx.Delete(&model.Product{}, []int64{id})

	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return result.Error
}
