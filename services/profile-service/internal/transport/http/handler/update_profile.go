package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sssoultrix/golive/services/profile-service/internal/usecase"
)

type UpdateProfileRequest struct {
	Name  string `json:"name" binding:"omitempty,min=2,max=100"`
	Email string `json:"email" binding:"omitempty,email"`
	Bio   string `json:"bio" binding:"omitempty,max=500"`
	Image string `json:"image" binding:"omitempty"`
}

type UpdateProfileResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	UpdatedAt string `json:"updated_at"`
}

type UpdateProfileHandler struct {
	uc *usecase.UpdateProfileUseCase
}

func NewUpdateProfileHandler(uc *usecase.UpdateProfileUseCase) *UpdateProfileHandler {
	return &UpdateProfileHandler{uc: uc}
}

func (h *UpdateProfileHandler) UpdateProfile(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	profile, err := h.uc.Execute(c.Request.Context(), id, req.Name, req.Email, req.Bio, req.Image)
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "update profile", slog.Any("err", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, UpdateProfileResponse{
		ID:        profile.ID.String(),
		Name:      profile.Name,
		Login:     profile.Login,
		Email:     profile.Email,
		Bio:       profile.Bio,
		Image:     profile.Image,
		UpdatedAt: profile.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

