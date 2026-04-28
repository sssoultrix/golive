package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/sssoultrix/golive/services/profile-service/internal/domain"
)

type GetProfileUseCase struct {
	profileRepo ProfileRepository
}

func NewGetProfileUseCase(profileRepo ProfileRepository) *GetProfileUseCase {
	return &GetProfileUseCase{
		profileRepo: profileRepo,
	}
}

func (g *GetProfileUseCase) ExecuteByID(ctx context.Context, id uuid.UUID) (domain.Profile, error) {
	return g.profileRepo.GetProfileByID(ctx, id)
}

func (g *GetProfileUseCase) ExecuteByLogin(ctx context.Context, login string) (domain.Profile, error) {
	return g.profileRepo.GetProfileByLogin(ctx, login)
}

