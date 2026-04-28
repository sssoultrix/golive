package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/sssoultrix/golive/services/profile-service/internal/domain"
)

type UpdateProfileUseCase struct {
	profileRepo ProfileRepository
}

func NewUpdateProfileUseCase(profileRepo ProfileRepository) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{
		profileRepo: profileRepo,
	}
}

func (u *UpdateProfileUseCase) Execute(ctx context.Context, id uuid.UUID, name, email, bio, image string) (domain.Profile, error) {
	profile, err := u.profileRepo.GetProfileByID(ctx, id)
	if err != nil {
		return domain.Profile{}, err
	}

	if name != "" {
		profile.Name = name
	}
	if email != "" {
		profile.Email = email
	}
	if bio != "" {
		profile.Bio = bio
	}
	if image != "" {
		profile.Image = image
	}

	if err := u.profileRepo.UpdateProfile(ctx, profile); err != nil {
		return domain.Profile{}, err
	}

	return profile, nil
}

