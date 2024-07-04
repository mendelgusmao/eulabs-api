package service

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/mendelgusmao/eulabs-api/domain"
	"github.com/mendelgusmao/eulabs-api/domain/dto"
	"github.com/mendelgusmao/eulabs-api/domain/model"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	FetchOne(context.Context, string) (*model.User, error)
	Create(context.Context, model.User) error
}

type UserService struct {
	repository UserRepository
}

func NewUserService(repository UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) Authorize(ctx context.Context, credentials dto.UserCredentials) (*dto.User, error) {
	user, err := s.repository.FetchOne(ctx, credentials.Username)

	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrCredentialsDontMatch
		}

		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(credentials.Username)); err != nil {
		return nil, domain.ErrCredentialsDontMatch
	}

	userRepresentation := &dto.User{}

	if err := copier.Copy(userRepresentation, user); err != nil {
		return nil, err
	}

	return userRepresentation, nil
}

func (s *UserService) Create(ctx context.Context, user dto.CreateUser) error {
	userModel := model.User{}

	if err := copier.Copy(&userModel, user); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	userModel.Password = hashedPassword

	return s.repository.Create(ctx, userModel)
}
