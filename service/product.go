package service

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/mendelgusmao/eulabs-api/domain/dto"
	"github.com/mendelgusmao/eulabs-api/domain/model"
)

type ProductRepository interface {
	FetchMany(context.Context, ...any) ([]model.Product, error)
	FetchOne(context.Context, int64) (*model.Product, error)
	Create(context.Context, model.Product) error
	Update(context.Context, model.Product) error
	Delete(context.Context, int64) error
}

type ProductService struct {
	repository ProductRepository
}

func NewProductService(repository ProductRepository) *ProductService {
	return &ProductService{
		repository: repository,
	}
}

func (s *ProductService) GetMany(ctx context.Context) ([]dto.Product, error) {
	products, err := s.repository.FetchMany(ctx)

	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return []dto.Product{}, nil
	}

	productsRepresentation := make([]dto.Product, 0)

	for _, product := range products {
		productRepresentation := dto.Product{}

		if err := copier.Copy(&productRepresentation, product); err != nil {
			return nil, err
		}

		productsRepresentation = append(productsRepresentation, productRepresentation)
	}

	return productsRepresentation, nil
}

func (s *ProductService) GetOne(ctx context.Context, id int64) (*dto.Product, error) {
	product, err := s.repository.FetchOne(ctx, id)

	if err != nil {
		return nil, err
	}

	productRepresentation := &dto.Product{}

	if err := copier.Copy(productRepresentation, product); err != nil {
		return nil, err
	}

	return productRepresentation, nil
}

func (s *ProductService) Create(ctx context.Context, product dto.BaseProduct) error {
	productModel := model.Product{}

	if err := copier.Copy(&productModel, product); err != nil {
		return err
	}

	return s.repository.Create(ctx, productModel)
}

func (s *ProductService) Update(ctx context.Context, product dto.Product) error {
	productModel := model.Product{}

	if err := copier.Copy(&productModel, product); err != nil {
		return err
	}

	return s.repository.Update(ctx, productModel)
}

func (s *ProductService) Delete(ctx context.Context, id int64) (err error) {
	return s.repository.Delete(ctx, id)
}