package domain

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Login     string    `json:"login,omitempty"`
	Email     string    `json:"email,omitempty"`
	Bio       string    `json:"bio,omitempty"`
	Image     string    `json:"image,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func normalizeAndValidateName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", fmt.Errorf("%w: empty", ErrInvalidName)
	}
	if len(name) < 2 || len(name) > 100 {
		return "", fmt.Errorf("%w: length must be between 2 and 100 characters", ErrInvalidName)
	}
	return name, nil
}

func normalizeAndValidateLogin(login string) (string, error) {
	login = strings.TrimSpace(login)
	if login == "" {
		return "", fmt.Errorf("%w: empty", ErrInvalidLogin)
	}
	if len(login) < 3 || len(login) > 32 {
		return "", fmt.Errorf("%w: length must be between 3 and 32 characters", ErrInvalidLogin)
	}
	// Allow only alphanumeric characters and underscores
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, login)
	if !matched {
		return "", fmt.Errorf("%w: only alphanumeric characters and underscores allowed", ErrInvalidLogin)
	}
	return login, nil
}

func normalizeAndValidateEmail(email string) (string, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return "", fmt.Errorf("%w: empty", ErrInvalidEmail)
	}
	// Basic email validation regex
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
	if !matched {
		return "", fmt.Errorf("%w: invalid format", ErrInvalidEmail)
	}
	return email, nil
}

func normalizeAndValidateBio(bio string) (string, error) {
	bio = strings.TrimSpace(bio)
	if len(bio) > 500 {
		return "", fmt.Errorf("%w: maximum 500 characters allowed", ErrInvalidBio)
	}
	return bio, nil
}

func NewProfile(name, login, email, bio, image string) (*Profile, error) {
	name, err := normalizeAndValidateName(name)
	if err != nil {
		return nil, err
	}

	login, err = normalizeAndValidateLogin(login)
	if err != nil {
		return nil, err
	}

	email, err = normalizeAndValidateEmail(email)
	if err != nil {
		return nil, err
	}

	bio, err = normalizeAndValidateBio(bio)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &Profile{
		ID:        uuid.New(),
		Name:      name,
		Login:     login,
		Email:     email,
		Bio:       bio,
		Image:     strings.TrimSpace(image),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
