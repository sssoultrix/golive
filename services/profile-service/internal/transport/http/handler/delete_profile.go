package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sssoultrix/golive/services/profile-service/internal/usecase"
)

type DeleteProfileHandler struct {
	uc *usecase.DeleteProfileUseCase
}

func NewDeleteProfileHandler(uc *usecase.DeleteProfileUseCase) *DeleteProfileHandler {
	return &DeleteProfileHandler{uc: uc}
}

func (h *DeleteProfileHandler) DeleteProfile(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.uc.Execute(c.Request.Context(), id); err != nil {
		slog.ErrorContext(c.Request.Context(), "delete profile", slog.Any("err", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete profile"})
		return
	}

	c.Status(http.StatusNoContent)
}

