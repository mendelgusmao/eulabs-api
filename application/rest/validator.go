package rest

import (
	"github.com/go-playground/validator"
)

type Validator struct {
	Validator *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		Validator: validator.New(),
	}
}

func (v *Validator) Validate(vv any) error {
	return v.Validator.Struct(vv)
}
