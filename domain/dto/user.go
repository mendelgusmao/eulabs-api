package dto

import "time"

type User struct {
	ID        int64
	Name      string
	Username  string
	Admin     bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateUser struct {
	ID        int64
	Name      string
	Username  string
	Password  string
	Admin     bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserCredentials struct {
	Username string
	Password string
}
