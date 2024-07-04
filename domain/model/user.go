package model

import "time"

type User struct {
	ID        int64 `gorm:"primarykey"`
	Name      string
	Username  string
	Password  []byte
	Admin     bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
