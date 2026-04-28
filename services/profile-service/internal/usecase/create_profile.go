package usecase

import (
	"context"

	"github.com/sssoultrix/golive/services/profile-service/internal/domain"
)

type CreateProfileUseCase struct {
	profileRepo ProfileRepository
}

func NewCreateProfileUseCase(profileRepo ProfileRepository) *CreateProfileUseCase {
	return &CreateProfileUseCase{
		profileRepo: profileRepo,
	}
}

func (c *CreateProfileUseCase) Execute(ctx context.Context, name, login, email, bio, image string) (domain.Profile, error) {
	profile, err := domain.NewProfile(name, login, email, bio, image)
	if err != nil {
		return domain.Profile{}, err
	}

	if err := c.profileRepo.CreateProfile(ctx, *profile); err != nil {
		return domain.Profile{}, err
	}

	return *profile, nil
}
