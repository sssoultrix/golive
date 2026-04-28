package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sssoultrix/golive/services/profile-service/internal/usecase"
)

type GetProfileResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type GetProfileHandler struct {
	uc *usecase.GetProfileUseCase
}

func NewGetProfileHandler(uc *usecase.GetProfileUseCase) *GetProfileHandler {
	return &GetProfileHandler{uc: uc}
}

func (h *GetProfileHandler) GetProfileByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	profile, err := h.uc.ExecuteByID(c.Request.Context(), id)
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "get profile by id", slog.Any("err", err))
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	c.JSON(http.StatusOK, GetProfileResponse{
		ID:        profile.ID.String(),
		Name:      profile.Name,
		Login:     profile.Login,
		Email:     profile.Email,
		Bio:       profile.Bio,
		Image:     profile.Image,
		CreatedAt: profile.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: profile.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

func (h *GetProfileHandler) GetProfileByLogin(c *gin.Context) {
	login := c.Param("login")

	profile, err := h.uc.ExecuteByLogin(c.Request.Context(), login)
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "get profile by login", slog.Any("err", err))
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	c.JSON(http.StatusOK, GetProfileResponse{
		ID:        profile.ID.String(),
		Name:      profile.Name,
		Login:     profile.Login,
		Email:     profile.Email,
		Bio:       profile.Bio,
		Image:     profile.Image,
		CreatedAt: profile.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: profile.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

