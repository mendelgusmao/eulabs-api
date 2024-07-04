package rest

import "fmt"

var (
	ErrInvalidIdType = Error(fmt.Errorf("invalid type for resource id"))
	ErrEmptyId       = Error(fmt.Errorf("resource id cant be empty"))
)

func Error(err error) map[string]string {
	return map[string]string{
		"error": err.Error(),
	}
}
