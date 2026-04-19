package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sssoultrix/golive/services/auth-service/internal/domain"
	"github.com/sssoultrix/golive/services/auth-service/internal/usecase"
)

type LoginRequest struct {
	Login    string `json:"login" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginResponse struct {
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginHandler struct {
	uc *usecase.LoginUser
}

func NewLoginHandler(uc *usecase.LoginUser) *LoginHandler {
	return &LoginHandler{uc: uc}
}

func (h *LoginHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	id, pair, err := h.uc.Execute(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		default:
			slog.ErrorContext(c.Request.Context(), "login user", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		}
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		UserID:       id.String(),
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
	})
}

