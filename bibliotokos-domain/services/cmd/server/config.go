package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"bibliotokos.domain/pkg/datastore"
	"bibliotokos.domain/pkg/user"
	"github.com/joho/godotenv"
)

type ServerConfig struct {
	UserDatabaseConfig user.Config `json:"user_database_config"`
	Port               int         `json:"port"`
	DB                 *datastore.DB
}

func LoadConfig() (ServerConfig, error) {
	if err := godotenv.Load(); err != nil {
		return ServerConfig{}, fmt.Errorf("loading .env: %w", err)
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		return ServerConfig{}, fmt.Errorf("DB_PATH is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return ServerConfig{}, fmt.Errorf("JWT_SECRET is required")
	}

	tokenExpiryMinutes := 60
	if v := os.Getenv("TOKEN_EXPIRY_MINUTES"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return ServerConfig{}, fmt.Errorf("invalid TOKEN_EXPIRY_MINUTES: %w", err)
		}
		tokenExpiryMinutes = n
	}

	port := 8080
	if v := os.Getenv("PORT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return ServerConfig{}, fmt.Errorf("invalid PORT: %w", err)
		}
		port = n
	}

	db, err := datastore.Open(dbPath)
	if err != nil {
		return ServerConfig{}, fmt.Errorf("opening database: %w", err)
	}

	return ServerConfig{
		UserDatabaseConfig: user.Config{
			DB:          db,
			JWTSecret:   jwtSecret,
			TokenExpiry: time.Duration(tokenExpiryMinutes) * time.Minute,
		},
		Port: port,
		DB:   db,
	}, nil
}
