package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"football-api/internal/core/config"
	"football-api/internal/core/middleware"
	"football-api/internal/domains/users/controllers"
	"football-api/internal/domains/users/repositories"
	"football-api/internal/domains/users/services"
	"football-api/internal/shared/constants"
)

func Register(api *gin.RouterGroup, db *gorm.DB, cfg config.Config) {
	repo := repositories.New(db)
	svc := services.New(repo)
	ctl := controllers.New(svc)

	group := api.Group("/users")
	group.Use(middleware.RequireAuth(cfg))
	{
		group.GET("", middleware.RequireRoles(string(constants.RoleAdmin)), ctl.List)
		group.POST("", middleware.RequireRoles(string(constants.RoleAdmin)), ctl.Create)
	}
}
