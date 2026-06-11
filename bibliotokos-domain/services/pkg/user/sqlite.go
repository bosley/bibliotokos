package user

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	_ "modernc.org/sqlite"
	"golang.org/x/crypto/bcrypt"
)
const schema = `
CREATE TABLE IF NOT EXISTS users (
	id            TEXT PRIMARY KEY,
	email         TEXT UNIQUE NOT NULL,
	password_hash TEXT NOT NULL,
	role          TEXT NOT NULL DEFAULT 'standard',
	verified      INTEGER NOT NULL DEFAULT 0,
	verified_at   DATETIME,
	created_at    DATETIME NOT NULL,
	updated_at    DATETIME NOT NULL,
	last_login    DATETIME,
	token_version INTEGER NOT NULL DEFAULT 0
);`

const migrateTokenVersion = `ALTER TABLE users ADD COLUMN token_version INTEGER NOT NULL DEFAULT 0;`

type sqliteController struct {
	writer *sql.DB
	reader *sql.DB
	cfg    Config
}

func newSQLiteController(cfg Config) (*sqliteController, error) {
	writer := cfg.DB.Writer()
	reader := cfg.DB.Reader()

	if _, err := writer.Exec(schema); err != nil {
		return nil, err
	}

	if _, err := writer.Exec(migrateTokenVersion); err != nil && !strings.Contains(err.Error(), "duplicate column name") {
		return nil, err
	}

	return &sqliteController{
		writer: writer,
		reader: reader,
		cfg:    cfg,
	}, nil
}

func isUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

func whereClause(identifier string) (string, string) {
	if isUUID(identifier) {
		return "id = ?", identifier
	}
	return "email = ?", identifier
}

func (c *sqliteController) CreateUser(email string, password string, role Role) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	id := uuid.NewString()
	now := time.Now().UTC()

	_, err = c.writer.Exec(
		`INSERT INTO users (id, email, password_hash, role, verified, created_at, updated_at)
		 VALUES (?, ?, ?, ?, 0, ?, ?)`,
		id, email, string(hash), string(role), now, now,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return "", ErrUserAlreadyExists
		}
		return "", err
	}

	return id, nil
}

func (c *sqliteController) DeleteUser(identifier string) error {
	col, val := whereClause(identifier)
	res, err := c.writer.Exec(
		fmt.Sprintf("DELETE FROM users WHERE %s", col),
		val,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (c *sqliteController) ResetPassword(identifier string, newPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	col, val := whereClause(identifier)
	res, err := c.writer.Exec(
		fmt.Sprintf("UPDATE users SET password_hash = ?, updated_at = ? WHERE %s", col),
		string(hash), time.Now().UTC(), val,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (c *sqliteController) SetRole(identifier string, role Role) error {
	col, val := whereClause(identifier)
	res, err := c.writer.Exec(
		fmt.Sprintf("UPDATE users SET role = ?, updated_at = ? WHERE %s", col),
		string(role), time.Now().UTC(), val,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (c *sqliteController) VerifyUser(identifier string) error {
	col, val := whereClause(identifier)
	now := time.Now().UTC()
	res, err := c.writer.Exec(
		fmt.Sprintf("UPDATE users SET verified = 1, verified_at = ?, updated_at = ? WHERE %s", col),
		now, now, val,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (c *sqliteController) tokenExpiry() time.Duration {
	if c.cfg.TokenExpiry <= 0 {
		return 24 * time.Hour
	}
	return c.cfg.TokenExpiry
}

func (c *sqliteController) Login(identifier string, password string) (UserRecord, string, error) {
	col, val := whereClause(identifier)

	row := c.reader.QueryRow(
		fmt.Sprintf(
			`SELECT id, email, password_hash, role, verified, verified_at, created_at, updated_at, last_login, token_version
			 FROM users WHERE %s`, col,
		),
		val,
	)

	var rec UserRecord
	var passwordHash string
	var verifiedAt, lastLogin sql.NullTime
	var tokenVersion int64

	err := row.Scan(
		&rec.ID, &rec.Email, &passwordHash, &rec.Role,
		&rec.Verified, &verifiedAt, &rec.CreatedAt, &rec.UpdatedAt, &lastLogin, &tokenVersion,
	)
	if err == sql.ErrNoRows {
		return UserRecord{}, "", ErrUserNotFound
	}
	if err != nil {
		return UserRecord{}, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return UserRecord{}, "", ErrInvalidPassword
	}

	if verifiedAt.Valid {
		rec.VerifiedAt = verifiedAt.Time
	}
	if lastLogin.Valid {
		rec.LastLogin = lastLogin.Time
	}

	now := time.Now().UTC()
	_, _ = c.writer.Exec(
		"UPDATE users SET last_login = ?, updated_at = ? WHERE id = ?",
		now, now, rec.ID,
	)
	rec.LastLogin = now

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   rec.ID,
		"email": rec.Email,
		"role":  string(rec.Role),
		"tkv":   tokenVersion,
		"exp":   now.Add(c.tokenExpiry()).Unix(),
	})

	signed, err := token.SignedString([]byte(c.cfg.JWTSecret))
	if err != nil {
		return UserRecord{}, "", err
	}

	return rec, signed, nil
}

func (c *sqliteController) Logout(identifier string) error {
	col, val := whereClause(identifier)
	res, err := c.writer.Exec(
		fmt.Sprintf("UPDATE users SET token_version = token_version + 1, updated_at = ? WHERE %s", col),
		time.Now().UTC(), val,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (c *sqliteController) ListUsers(offset, limit int) ([]UserRecord, error) {
	rows, err := c.reader.Query(
		`SELECT id, email, role, verified, verified_at, created_at, updated_at, last_login
		 FROM users ORDER BY created_at ASC LIMIT ? OFFSET ?`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []UserRecord
	for rows.Next() {
		var rec UserRecord
		var verifiedAt, lastLogin sql.NullTime
		if err := rows.Scan(
			&rec.ID, &rec.Email, &rec.Role, &rec.Verified,
			&verifiedAt, &rec.CreatedAt, &rec.UpdatedAt, &lastLogin,
		); err != nil {
			return nil, err
		}
		if verifiedAt.Valid {
			rec.VerifiedAt = verifiedAt.Time
		}
		if lastLogin.Valid {
			rec.LastLogin = lastLogin.Time
		}
		users = append(users, rec)
	}
	return users, rows.Err()
}

func (c *sqliteController) ValidateToken(tokenString string) (UserRecord, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(c.cfg.JWTSecret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return UserRecord{}, ErrTokenExpired
		}
		return UserRecord{}, ErrInvalidToken
	}
	if !token.Valid {
		return UserRecord{}, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return UserRecord{}, ErrInvalidToken
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return UserRecord{}, ErrInvalidToken
	}

	tkvFloat, ok := claims["tkv"].(float64)
	if !ok {
		return UserRecord{}, ErrInvalidToken
	}
	claimedVersion := int64(tkvFloat)

	row := c.reader.QueryRow(
		`SELECT id, email, password_hash, role, verified, verified_at, created_at, updated_at, last_login, token_version
		 FROM users WHERE id = ?`,
		sub,
	)

	var rec UserRecord
	var passwordHash string
	var verifiedAt, lastLogin sql.NullTime
	var tokenVersion int64

	if err := row.Scan(
		&rec.ID, &rec.Email, &passwordHash, &rec.Role,
		&rec.Verified, &verifiedAt, &rec.CreatedAt, &rec.UpdatedAt, &lastLogin, &tokenVersion,
	); err == sql.ErrNoRows {
		return UserRecord{}, ErrUserNotFound
	} else if err != nil {
		return UserRecord{}, err
	}

	if claimedVersion != tokenVersion {
		return UserRecord{}, ErrInvalidToken
	}

	if verifiedAt.Valid {
		rec.VerifiedAt = verifiedAt.Time
	}
	if lastLogin.Valid {
		rec.LastLogin = lastLogin.Time
	}

	return rec, nil
}
