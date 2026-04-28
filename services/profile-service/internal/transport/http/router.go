package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sssoultrix/golive/services/profile-service/internal/transport/http/handler"
	"github.com/sssoultrix/golive/services/profile-service/internal/usecase"
)

type Router struct {
	engine *gin.Engine
}

func NewRouter() *Router {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())

	return &Router{
		engine: engine,
	}
}

func (r *Router) SetupRoutes(
	createUC *usecase.CreateProfileUseCase,
	getUC *usecase.GetProfileUseCase,
	updateUC *usecase.UpdateProfileUseCase,
	deleteUC *usecase.DeleteProfileUseCase,
) {
	// Health check endpoint
	r.engine.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Profile endpoints
	createHandler := handler.NewCreateProfileHandler(createUC)
	r.engine.POST("/profiles", createHandler.CreateProfile)

	getHandler := handler.NewGetProfileHandler(getUC)
	r.engine.GET("/profiles/:id", getHandler.GetProfileByID)
	r.engine.GET("/profiles/login/:login", getHandler.GetProfileByLogin)

	updateHandler := handler.NewUpdateProfileHandler(updateUC)
	r.engine.PUT("/profiles/:id", updateHandler.UpdateProfile)

	deleteHandler := handler.NewDeleteProfileHandler(deleteUC)
	r.engine.DELETE("/profiles/:id", deleteHandler.DeleteProfile)
}

func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}
