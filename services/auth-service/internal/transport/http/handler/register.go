package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sssoultrix/golive/services/auth-service/internal/domain"
	"github.com/sssoultrix/golive/services/auth-service/internal/usecase"
)

const postgresUniqueViolation = "23505"

type RegisterRequest struct {
	Login    string `json:"login" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=8"`
}

type RegisterResponse struct {
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterHandler struct {
	uc *usecase.RegisterUser
}

func NewRegisterHandler(uc *usecase.RegisterUser) *RegisterHandler {
	return &RegisterHandler{uc: uc}
}

func (h *RegisterHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	id, pair, err := h.uc.Execute(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidLogin):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		case isPostgresUniqueViolation(err):
			c.JSON(http.StatusConflict, gin.H{"error": "login already taken"})
		default:
			slog.ErrorContext(c.Request.Context(), "register user", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		}
		return
	}

	c.JSON(http.StatusCreated, RegisterResponse{
		UserID:       id.String(),
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
	})
}

func isPostgresUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == postgresUniqueViolation
}
