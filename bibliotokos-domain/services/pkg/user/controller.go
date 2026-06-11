package user

import (
	"time"

	"bibliotokos.domain/pkg/datastore"
)

type Config struct {
	DB          *datastore.DB
	JWTSecret   string
	TokenExpiry time.Duration
}

type UserController interface {
	CreateUser(email string, password string, role Role) (string, error)
	DeleteUser(identifier string) error
	ResetPassword(identifier string, newPassword string) error
	SetRole(identifier string, role Role) error
	VerifyUser(identifier string) error
	Login(identifier string, password string) (record UserRecord, jwt string, err error)
	Logout(identifier string) error
	ValidateToken(tokenString string) (UserRecord, error)
	ListUsers(offset, limit int) ([]UserRecord, error)
}

func FromDB(cfg Config) (UserController, error) {
	return newSQLiteController(cfg)
}
