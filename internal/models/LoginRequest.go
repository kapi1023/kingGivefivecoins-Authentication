package models

type LoginRequest struct {
	Email    string `json:"email" example:"email@example.com" validate:"required,email"`
	Password string `json:"password" example:"password" validate:"required"`
}
