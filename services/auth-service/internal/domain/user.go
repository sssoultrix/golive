package domain

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

func isLoginASCIIRune(r rune) bool {
	switch {
	case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
		return true
	case r == '_':
		return true
	default:
		return false
	}
}

type User struct {
	ID           uuid.UUID `json:"id"`
	Login        string    `json:"login"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func normalizeAndValidateLogin(login string) (string, error) {
	login = strings.TrimSpace(login)
	if login == "" {
		return "", fmt.Errorf("%w: empty", ErrInvalidLogin)
	}
	n := utf8.RuneCountInString(login)
	if n < 3 || n > 32 {
		return "", fmt.Errorf("%w: length", ErrInvalidLogin)
	}
	for _, r := range login {
		if !isLoginASCIIRune(r) {
			return "", fmt.Errorf("%w: charset", ErrInvalidLogin)
		}
	}
	return login, nil
}

func NewUser(login, passwordHash string) (*User, error) {
	login, err := normalizeAndValidateLogin(login)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	return &User{
		ID:           uuid.New(),
		Login:        login,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}
