package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"football-api/internal/core/config"
	"football-api/internal/core/middleware"
	"football-api/internal/domains/auth/controllers"
	"football-api/internal/domains/auth/repositories"
	"football-api/internal/domains/auth/services"
	userroutes "football-api/internal/domains/users/routes"
)

func Register(api *gin.RouterGroup, db *gorm.DB, cfg config.Config) {
	repo := repositories.New(db)
	svc := services.New(repo, cfg)
	ctl := controllers.New(svc)

	auth := api.Group("/auth")
	{
		auth.POST("/login", middleware.LoginRateLimit(cfg), ctl.Login)
		auth.POST("/refresh", ctl.Refresh)
	}

	userroutes.Register(api, db, cfg)
}
