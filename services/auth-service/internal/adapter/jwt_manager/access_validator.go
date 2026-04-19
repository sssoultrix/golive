package jwt_manager

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sssoultrix/golive/services/auth-service/internal/domain"
)

type AccessValidator struct {
	secret []byte
}

func NewAccessValidator(secret []byte) *AccessValidator {
	return &AccessValidator{secret: secret}
}

func (v *AccessValidator) Validate(tokenString string) (jwt.RegisteredClaims, error) {
	var claims jwt.RegisteredClaims

	parsed, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return v.secret, nil
	})
	if err != nil || parsed == nil || !parsed.Valid {
		return jwt.RegisteredClaims{}, domain.ErrInvalidAccessToken
	}

	return claims, nil
}

