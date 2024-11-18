package models

import "time"

type User struct {
	ID            int
	Email         string
	PasswordHash  *string
	OAuthProvider *string
	OAuthID       *string
	CreatedAt     time.Time
}
