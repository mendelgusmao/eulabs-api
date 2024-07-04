package domain

import "fmt"

var (
	ErrNotFound             = fmt.Errorf("record not found")
	ErrCredentialsDontMatch = fmt.Errorf("credentials dont match")
)
