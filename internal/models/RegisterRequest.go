package models

type RegisterRequest struct {
	Email    string `json:"email" example:"user@example.com" validate:"required,email"`
	Password string `json:"password" example:"securepassword" validate:"required,min=8"`
}
