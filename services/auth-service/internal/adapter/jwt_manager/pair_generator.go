package jwt_manager

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sssoultrix/golive/services/auth-service/internal/usecase"
)

type PairGenerator struct {
	secret       []byte
	accessTTL    time.Duration
	refreshBytes int
}

func NewPairGenerator(secret []byte, accessTTL time.Duration, refreshBytes int) *PairGenerator {
	if refreshBytes < 16 {
		refreshBytes = 32
	}
	return &PairGenerator{
		secret,
		accessTTL,
		refreshBytes,
	}
}

func (p *PairGenerator) GeneratePair(userID uuid.UUID) (usecase.TokenPair, error) {
	now := time.Now()

	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(p.accessTTL)),
		ID:        uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := token.SignedString(p.secret)
	if err != nil {
		return usecase.TokenPair{}, err
	}

	buf := make([]byte, p.refreshBytes)
	if _, err := rand.Read(buf); err != nil {
		return usecase.TokenPair{}, err
	}

	refresh := base64.RawURLEncoding.EncodeToString(buf)

	return usecase.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refresh,
	}, nil
}
