package models

import "time"

type User struct {
	UserId    int     `json:"user_id"`
	Username  string  `json:"username" validate:"required" msg:""`
	Password  string  `json:"password" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserResponse struct {
	UserId    int     `json:"user_id"`
	Username  string  `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type UserResponsePagination struct {
	Users  []UserResponse `json:"users"`
	Total  int            `json:"total"`
	Page   int            `json:"page"`
	Limit  int            `json:"limit"`
}
