package user

import (
	"errors"
	"os"
	"testing"

	"bibliotokos.domain/pkg/datastore"
)

var testCtrl UserController

func TestMain(m *testing.M) {
	f, err := os.CreateTemp("", "user_test_*.db")
	if err != nil {
		panic(err)
	}
	dbPath := f.Name()
	f.Close()
	defer os.Remove(dbPath)

	ds, err := datastore.Open(dbPath)
	if err != nil {
		panic(err)
	}
	defer ds.Close()

	cfg := Config{
		DB:        ds,
		JWTSecret: "test-secret",
	}
	testCtrl, err = FromDB(cfg)
	if err != nil {
		panic(err)
	}

	code := m.Run()

	raw := testCtrl.(*sqliteController)
	raw.writer.Exec("DELETE FROM users WHERE email LIKE '%_test'")

	os.Exit(code)
}

func TestCreateUser(t *testing.T) {
	id, err := testCtrl.CreateUser("alice_test", "password123", RoleStandard)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id == "" {
		t.Fatal("expected non-empty user id")
	}

	_, err = testCtrl.CreateUser("alice_test", "password123", RoleStandard)
	if !errors.Is(err, ErrUserAlreadyExists) {
		t.Fatalf("expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestDeleteUser(t *testing.T) {
	id, err := testCtrl.CreateUser("delete_test", "password123", RoleStandard)
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := testCtrl.DeleteUser("delete_test"); err != nil {
		t.Fatalf("delete by email: %v", err)
	}

	id2, err := testCtrl.CreateUser("delete2_test", "password123", RoleStandard)
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := testCtrl.DeleteUser(id2); err != nil {
		t.Fatalf("delete by id: %v", err)
	}

	_ = id

	if err := testCtrl.DeleteUser("nonexistent_test"); !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestResetPassword(t *testing.T) {
	_, err := testCtrl.CreateUser("reset_test", "oldpass", RoleStandard)
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := testCtrl.ResetPassword("reset_test", "newpass"); err != nil {
		t.Fatalf("reset password: %v", err)
	}

	_, _, err = testCtrl.Login("reset_test", "oldpass")
	if !errors.Is(err, ErrInvalidPassword) {
		t.Fatalf("expected old password to be invalid, got %v", err)
	}

	if err := testCtrl.VerifyUser("reset_test"); err != nil {
		t.Fatalf("verify: %v", err)
	}

	_, _, err = testCtrl.Login("reset_test", "newpass")
	if err != nil {
		t.Fatalf("expected login with new password to succeed, got %v", err)
	}

	if err := testCtrl.ResetPassword("ghost_test", "x"); !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestSetRole(t *testing.T) {
	_, err := testCtrl.CreateUser("role_test", "pass", RoleStandard)
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := testCtrl.SetRole("role_test", RoleAdmin); err != nil {
		t.Fatalf("set role: %v", err)
	}

	if err := testCtrl.VerifyUser("role_test"); err != nil {
		t.Fatalf("verify: %v", err)
	}

	rec, _, err := testCtrl.Login("role_test", "pass")
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if rec.Role != RoleAdmin {
		t.Fatalf("expected role admin, got %v", rec.Role)
	}

	if err := testCtrl.SetRole("ghost_test", RoleAdmin); !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestVerifyUser(t *testing.T) {
	_, err := testCtrl.CreateUser("verify_test", "pass", RoleStandard)
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := testCtrl.VerifyUser("verify_test"); err != nil {
		t.Fatalf("verify: %v", err)
	}

	if err := testCtrl.VerifyUser("ghost_test"); !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestLogin(t *testing.T) {
	_, err := testCtrl.CreateUser("login_test", "secret", RoleStandard)
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	if err := testCtrl.VerifyUser("login_test"); err != nil {
		t.Fatalf("verify: %v", err)
	}

	rec, jwtStr, err := testCtrl.Login("login_test", "secret")
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if rec.Email != "login_test" {
		t.Fatalf("expected email login_test, got %s", rec.Email)
	}
	if jwtStr == "" {
		t.Fatal("expected non-empty JWT")
	}
	if rec.LastLogin.IsZero() {
		t.Fatal("expected last_login to be set")
	}

	_, _, err = testCtrl.Login("login_test", "wrongpass")
	if !errors.Is(err, ErrInvalidPassword) {
		t.Fatalf("expected ErrInvalidPassword, got %v", err)
	}

	_, _, err = testCtrl.Login("ghost_test", "pass")
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}

	idRec, _, _ := testCtrl.Login("login_test", "secret")
	recByID, jwtByID, err := testCtrl.Login(idRec.ID, "secret")
	if err != nil {
		t.Fatalf("login by id: %v", err)
	}
	if recByID.ID != idRec.ID {
		t.Fatal("id mismatch when logging in by uuid")
	}
	if jwtByID == "" {
		t.Fatal("expected non-empty JWT when logging in by uuid")
	}
}

func TestValidateToken(t *testing.T) {
	_, err := testCtrl.CreateUser("validatetoken_test", "pass", RoleStandard)
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	if err := testCtrl.VerifyUser("validatetoken_test"); err != nil {
		t.Fatalf("verify: %v", err)
	}

	_, token, err := testCtrl.Login("validatetoken_test", "pass")
	if err != nil {
		t.Fatalf("login: %v", err)
	}

	rec, err := testCtrl.ValidateToken(token)
	if err != nil {
		t.Fatalf("expected valid token, got %v", err)
	}
	if rec.Email != "validatetoken_test" {
		t.Fatalf("expected email validatetoken_test, got %s", rec.Email)
	}

	_, err = testCtrl.ValidateToken("not.a.valid.token")
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken for garbage input, got %v", err)
	}

	_, err = testCtrl.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyMTIzIiwiZW1haWwiOiJ4QHguY29tIiwicm9sZSI6InN0YW5kYXJkIiwiZXhwIjoxfQ.invalidsig")
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken for bad signature, got %v", err)
	}

	_, err = testCtrl.ValidateToken("")
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken for empty string, got %v", err)
	}
}

func TestLogout(t *testing.T) {
	_, err := testCtrl.CreateUser("logout_test", "pass", RoleStandard)
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	if err := testCtrl.VerifyUser("logout_test"); err != nil {
		t.Fatalf("verify: %v", err)
	}

	_, token, err := testCtrl.Login("logout_test", "pass")
	if err != nil {
		t.Fatalf("login: %v", err)
	}

	if _, err := testCtrl.ValidateToken(token); err != nil {
		t.Fatalf("token should be valid before logout, got %v", err)
	}

	if err := testCtrl.Logout("logout_test"); err != nil {
		t.Fatalf("logout: %v", err)
	}

	_, err = testCtrl.ValidateToken(token)
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken after logout, got %v", err)
	}

	_, newToken, err := testCtrl.Login("logout_test", "pass")
	if err != nil {
		t.Fatalf("re-login after logout: %v", err)
	}
	if _, err := testCtrl.ValidateToken(newToken); err != nil {
		t.Fatalf("new token should be valid after re-login, got %v", err)
	}

	if err := testCtrl.Logout("ghost_test"); !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
