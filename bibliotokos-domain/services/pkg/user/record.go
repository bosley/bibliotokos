package user

import "time"

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleStandard Role = "standard"
)

type UserRecord struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Role       Role      `json:"role"`
	Verified   bool      `json:"verified"`
	VerifiedAt time.Time `json:"verified_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	LastLogin  time.Time `json:"last_login"` // set to never
}
