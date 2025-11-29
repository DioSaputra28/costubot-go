package models

import "time"

type User struct {
	UserId    int     `json:"user_id"`
	Username  string  `json:"username" validate:"required"`
	Password  string  `json:"password" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
