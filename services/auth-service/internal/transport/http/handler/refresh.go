package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sssoultrix/golive/services/auth-service/internal/domain"
	"github.com/sssoultrix/golive/services/auth-service/internal/usecase"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshHandler struct {
	uc *usecase.RefreshTokens
}

func NewRefreshHandler(uc *usecase.RefreshTokens) *RefreshHandler {
	return &RefreshHandler{uc: uc}
}

func (h *RefreshHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	pair, err := h.uc.Execute(c.Request.Context(), req.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidRefreshToken):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		default:
			slog.ErrorContext(c.Request.Context(), "refresh tokens", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "refresh failed"})
		}
		return
	}

	c.JSON(http.StatusOK, RefreshResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
	})
}

