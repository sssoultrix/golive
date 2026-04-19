package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sssoultrix/golive/services/auth-service/internal/domain"
	"github.com/sssoultrix/golive/services/auth-service/internal/usecase"
)

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutHandler struct {
	uc *usecase.Logout
}

func NewLogoutHandler(uc *usecase.Logout) *LogoutHandler {
	return &LogoutHandler{uc: uc}
}

func (h *LogoutHandler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	if err := h.uc.Execute(c.Request.Context(), req.RefreshToken); err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidRefreshToken):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		default:
			slog.ErrorContext(c.Request.Context(), "logout", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "logout failed"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

