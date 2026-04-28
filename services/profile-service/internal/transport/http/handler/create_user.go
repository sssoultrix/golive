package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sssoultrix/golive/services/profile-service/internal/domain"
	"github.com/sssoultrix/golive/services/profile-service/internal/usecase"
)

const postgresUniqueViolation = "23505"

type CreateProfileRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=100"`
	Login string `json:"login" binding:"required,min=3,max=32"`
	Email string `json:"email" binding:"required,email"`
	Bio   string `json:"bio" binding:"omitempty,max=500"`
	Image string `json:"image" binding:"omitempty"`
}

type CreateProfileResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Login string `json:"login"`
	Email string `json:"email"`
}

type CreateProfileHandler struct {
	uc *usecase.CreateProfileUseCase
}

func NewCreateProfileHandler(uc *usecase.CreateProfileUseCase) *CreateProfileHandler {
	return &CreateProfileHandler{uc: uc}
}

func (h *CreateProfileHandler) CreateProfile(c *gin.Context) {
	var req CreateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	profile, err := h.uc.Execute(c.Request.Context(), req.Name, req.Login, req.Email, req.Bio, req.Image)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidName), errors.Is(err, domain.ErrInvalidLogin), errors.Is(err, domain.ErrInvalidEmail), errors.Is(err, domain.ErrInvalidBio):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		case isPostgresUniqueViolation(err):
			c.JSON(http.StatusConflict, gin.H{"error": "login or email already taken"})
		default:
			slog.ErrorContext(c.Request.Context(), "create profile", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create profile"})
		}
		return
	}

	c.JSON(http.StatusCreated, CreateProfileResponse{
		ID:    profile.ID.String(),
		Name:  profile.Name,
		Login: profile.Login,
		Email: profile.Email,
	})
}

func isPostgresUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == postgresUniqueViolation
}

