package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"bibliotokos.domain/pkg/user"
)

type server struct {
	mu         sync.Mutex
	controller user.UserController
}

func (s *server) writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func (s *server) writeError(w http.ResponseWriter, status int, msg string) {
	s.writeJSON(w, status, map[string]string{"error": msg})
}

func (s *server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Email == "" || req.Password == "" {
		s.writeError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	s.mu.Lock()
	_, token, err := s.controller.Login(req.Email, req.Password)
	s.mu.Unlock()

	if err != nil {
		switch {
		case errors.Is(err, user.ErrUserNotFound), errors.Is(err, user.ErrInvalidPassword):
			s.writeError(w, http.StatusUnauthorized, "invalid credentials")
		case errors.Is(err, user.ErrUserNotVerified):
			s.writeError(w, http.StatusForbidden, "user not verified")
		default:
			s.writeError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	s.writeJSON(w, http.StatusOK, LoginResponse{Token: token})
}

func (s *server) handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Token == "" {
		s.writeError(w, http.StatusBadRequest, "token is required")
		return
	}

	s.mu.Lock()
	record, err := s.controller.ValidateToken(req.Token)
	s.mu.Unlock()

	if err != nil {
		switch {
		case errors.Is(err, user.ErrInvalidToken):
			s.writeError(w, http.StatusUnauthorized, "invalid token")
		case errors.Is(err, user.ErrTokenExpired):
			s.writeError(w, http.StatusUnauthorized, "token expired")
		default:
			s.writeError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	s.mu.Lock()
	err = s.controller.Logout(record.ID)
	s.mu.Unlock()

	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *server) handleValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Token == "" {
		s.writeError(w, http.StatusBadRequest, "token is required")
		return
	}

	s.mu.Lock()
	record, err := s.controller.ValidateToken(req.Token)
	s.mu.Unlock()

	if err != nil {
		switch {
		case errors.Is(err, user.ErrInvalidToken):
			s.writeError(w, http.StatusUnauthorized, "invalid token")
		case errors.Is(err, user.ErrTokenExpired):
			s.writeError(w, http.StatusUnauthorized, "token expired")
		default:
			s.writeError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	s.writeJSON(w, http.StatusOK, UserDTO{
		Email:    record.Email,
		Role:     string(record.Role),
		Verified: record.Verified,
		IsAdmin:  record.IsAdmin(),
	})
}

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	defer cfg.DB.Close()

	userController, err := user.FromDB(cfg.UserDatabaseConfig)
	if err != nil {
		log.Fatalf("failed to create user controller: %v", err)
	}

	srv := &server{controller: userController}

	mux := http.NewServeMux()
	mux.HandleFunc("/login", srv.handleLogin)
	mux.HandleFunc("/logout", srv.handleLogout)
	mux.HandleFunc("/validate", srv.handleValidate)

	addr := fmt.Sprintf("127.0.0.1:%d", cfg.Port)
	log.Printf("listening on %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
