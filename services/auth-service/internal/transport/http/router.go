package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sssoultrix/golive/services/auth-service/internal/adapter/jwt_manager"
	"github.com/sssoultrix/golive/services/auth-service/internal/transport/http/handler"
	"github.com/sssoultrix/golive/services/auth-service/internal/usecase"
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
	registerUC *usecase.RegisterUser,
	loginUC *usecase.LoginUser,
	refreshUC *usecase.RefreshTokens,
	logoutUC *usecase.Logout,
	accessValidator *jwt_manager.AccessValidator,
) {
	// Health check endpoint
	r.engine.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Auth endpoints
	registerHandler := handler.NewRegisterHandler(registerUC)
	r.engine.POST("/register", registerHandler.Register)

	loginHandler := handler.NewLoginHandler(loginUC)
	r.engine.POST("/login", loginHandler.Login)

	refreshHandler := handler.NewRefreshHandler(refreshUC)
	r.engine.POST("/refresh", refreshHandler.Refresh)

	logoutHandler := handler.NewLogoutHandler(logoutUC)
	r.engine.POST("/logout", logoutHandler.Logout)

	validateHandler := handler.NewValidateHandler(accessValidator)
	r.engine.POST("/validate", validateHandler.Validate)
}

func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}
