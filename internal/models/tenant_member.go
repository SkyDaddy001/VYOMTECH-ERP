package models

import "time"

// TenantMember represents a user's membership in a tenant
type TenantMember struct {
	ID        string    `db:"id" json:"id"`
	TenantID  string    `db:"tenant_id" json:"tenant_id"`
	UserID    int64     `db:"user_id" json:"user_id"`
	Email     string    `db:"email" json:"email"`
	Role      string    `db:"role" json:"role"` // admin, member, viewer
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
