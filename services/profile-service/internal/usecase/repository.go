package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/sssoultrix/golive/services/profile-service/internal/domain"
)

type ProfileRepository interface {
	CreateProfile(ctx context.Context, profile domain.Profile) error
	GetProfileByID(ctx context.Context, id uuid.UUID) (domain.Profile, error)
	GetProfileByLogin(ctx context.Context, login string) (domain.Profile, error)
	UpdateProfile(ctx context.Context, profile domain.Profile) error
	DeleteProfile(ctx context.Context, id uuid.UUID) error
}
