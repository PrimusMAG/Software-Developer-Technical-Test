package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"football-api/internal/core/config"
	"football-api/internal/core/middleware"
	"football-api/internal/domains/matches/controllers"
	"football-api/internal/domains/matches/repositories"
	"football-api/internal/domains/matches/services"
	"football-api/internal/shared/constants"
)

func Register(api *gin.RouterGroup, db *gorm.DB, cfg config.Config) {
	repo := repositories.New(db)
	svc := services.New(repo)
	ctl := controllers.New(svc)

	group := api.Group("/matches")
	group.Use(middleware.RequireAuth(cfg))
	{
		group.GET("", middleware.RequireRoles(string(constants.RoleAdmin), string(constants.RoleStaff), string(constants.RoleViewer)), ctl.List)
		group.GET("/:id", middleware.RequireRoles(string(constants.RoleAdmin), string(constants.RoleStaff), string(constants.RoleViewer)), ctl.GetByID)
		group.POST("", middleware.RequireRoles(string(constants.RoleAdmin)), ctl.Create)
		group.DELETE("/:id", middleware.RequireRoles(string(constants.RoleAdmin)), ctl.Delete)
		group.POST("/:id/result", middleware.RequireRoles(string(constants.RoleAdmin)), ctl.SubmitResult)
		group.POST("/:id/rollback", middleware.RequireRoles(string(constants.RoleAdmin)), ctl.RollbackResult)
	}
}
