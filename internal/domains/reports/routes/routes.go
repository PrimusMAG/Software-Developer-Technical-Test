package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"football-api/internal/core/config"
	"football-api/internal/core/middleware"
	"football-api/internal/domains/reports/controllers"
	"football-api/internal/domains/reports/repositories"
	"football-api/internal/domains/reports/services"
	"football-api/internal/shared/constants"
)

func Register(api *gin.RouterGroup, db *gorm.DB, cfg config.Config) {
	repo := repositories.New(db)
	svc := services.New(repo)
	ctl := controllers.New(svc)

	group := api.Group("/reports")
	group.Use(middleware.RequireAuth(cfg))
	{
		group.GET("/matches", middleware.RequireRoles(string(constants.RoleAdmin), string(constants.RoleStaff), string(constants.RoleViewer)), ctl.List)
		group.POST("/revalidate", middleware.RequireRoles(string(constants.RoleAdmin)), ctl.Revalidate)
	}
}
