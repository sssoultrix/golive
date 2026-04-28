package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sssoultrix/golive/services/profile-service/internal/domain"
)

type ProfileRepository struct {
	db *pgxpool.Pool
}

func NewProfileRepository(db *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (p *ProfileRepository) CreateProfile(ctx context.Context, profile domain.Profile) error {
	query := `INSERT INTO profiles (id, name, login, email, bio, image, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := p.db.Exec(ctx, query, profile.ID, profile.Name, profile.Login, profile.Email, profile.Bio, profile.Image, profile.CreatedAt, profile.UpdatedAt)
	return err
}

func (p *ProfileRepository) GetProfileByID(ctx context.Context, id uuid.UUID) (domain.Profile, error) {
	query := `SELECT id, name, login, email, bio, image, created_at, updated_at 
	          FROM profiles WHERE id = $1`

	var profile domain.Profile
	err := p.db.QueryRow(ctx, query, id).Scan(
		&profile.ID,
		&profile.Name,
		&profile.Login,
		&profile.Email,
		&profile.Bio,
		&profile.Image,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err != nil {
		return domain.Profile{}, err
	}
	return profile, nil
}

func (p *ProfileRepository) GetProfileByLogin(ctx context.Context, login string) (domain.Profile, error) {
	query := `SELECT id, name, login, email, bio, image, created_at, updated_at 
	          FROM profiles WHERE login = $1`

	var profile domain.Profile
	err := p.db.QueryRow(ctx, query, login).Scan(
		&profile.ID,
		&profile.Name,
		&profile.Login,
		&profile.Email,
		&profile.Bio,
		&profile.Image,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err != nil {
		return domain.Profile{}, err
	}
	return profile, nil
}

func (p *ProfileRepository) UpdateProfile(ctx context.Context, profile domain.Profile) error {
	query := `UPDATE profiles 
	          SET name = $2, email = $3, bio = $4, image = $5, updated_at = $6 
	          WHERE id = $1`

	profile.UpdatedAt = time.Now().UTC()
	_, err := p.db.Exec(ctx, query, profile.ID, profile.Name, profile.Email, profile.Bio, profile.Image, profile.UpdatedAt)
	return err
}

func (p *ProfileRepository) DeleteProfile(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM profiles WHERE id = $1`
	_, err := p.db.Exec(ctx, query, id)
	return err
}
