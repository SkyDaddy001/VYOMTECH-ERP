package models

import (
	"time"
)

type User struct {
	ID              int       `json:"id" db:"id"`
	Email           string    `json:"email" db:"email"`
	Password        string    `json:"-" db:"password_hash"`
	Role            string    `json:"role" db:"role"`
	TenantID        string    `json:"tenant_id" db:"tenant_id"`
	CurrentTenantID string    `json:"current_tenant_id" db:"current_tenant_id"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type PasswordResetToken struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
