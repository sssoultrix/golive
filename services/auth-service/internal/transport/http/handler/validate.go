package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sssoultrix/golive/services/auth-service/internal/domain"
)

type AccessTokenValidator interface {
	Validate(tokenString string) (jwt.RegisteredClaims, error)
}

type ValidateRequest struct {
	AccessToken string `json:"access_token" binding:"required"`
}

type ValidateResponse struct {
	UserID string `json:"user_id"`
	JTI    string `json:"jti,omitempty"`
	Exp    int64  `json:"exp,omitempty"`
}

type ValidateHandler struct {
	v AccessTokenValidator
}

func NewValidateHandler(v AccessTokenValidator) *ValidateHandler {
	return &ValidateHandler{v: v}
}

func (h *ValidateHandler) Validate(c *gin.Context) {
	accessToken := tokenFromAuthorizationHeader(c.GetHeader("Authorization"))
	if accessToken == "" {
		var req ValidateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}
		accessToken = req.AccessToken
	}

	claims, err := h.v.Validate(accessToken)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidAccessToken):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
		default:
			slog.ErrorContext(c.Request.Context(), "validate access token", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "validation failed"})
		}
		return
	}

	resp := ValidateResponse{
		UserID: claims.Subject,
		JTI:    claims.ID,
	}
	if claims.ExpiresAt != nil {
		resp.Exp = claims.ExpiresAt.Time.UTC().Unix()
	}

	c.JSON(http.StatusOK, resp)
}

func tokenFromAuthorizationHeader(authHeader string) string {
	authHeader = strings.TrimSpace(authHeader)
	if authHeader == "" {
		return ""
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(authHeader, prefix))
}

