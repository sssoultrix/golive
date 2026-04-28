package usecase

import (
	"context"

	"github.com/google/uuid"
)

type DeleteProfileUseCase struct {
	profileRepo ProfileRepository
}

func NewDeleteProfileUseCase(profileRepo ProfileRepository) *DeleteProfileUseCase {
	return &DeleteProfileUseCase{
		profileRepo: profileRepo,
	}
}

func (d *DeleteProfileUseCase) Execute(ctx context.Context, id uuid.UUID) error {
	return d.profileRepo.DeleteProfile(ctx, id)
}

