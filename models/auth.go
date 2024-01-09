package models

import "gorm.io/gorm"

// Register
type User struct {
	gorm.Model
	Email    string `db:"email"`
	Password string `db:"email"`
	Name     string `db:"name"`
	Age      int    `db:"age"`
	Phone    string `db:"phone"`
	Rank     string `db:"rank"`
	Role     string `db:"role"`
	Status   string `db:"status"`
}

type UserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Age      int    `json:"age"`
	Phone    string `json:"phone"`
	Rank     string `json:"rank"`
}

type UserResponse struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Rank   string `json:"rank"`
	Status string `json:"status"`
	Role   string `json:"role"`
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Redis struct {
	Key   string
	Value string
}

// Login
type LoginResponse struct {
	Status        string
	Access_token  string
	Refresh_token string
}
