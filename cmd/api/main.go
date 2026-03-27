package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"football-api/internal/core/config"
	"football-api/internal/core/database"
	"football-api/internal/core/middleware"
	"football-api/internal/domains/auth/routes"
	matchroutes "football-api/internal/domains/matches/routes"
	playerroutes "football-api/internal/domains/players/routes"
	reportroutes "football-api/internal/domains/reports/routes"
	teamroutes "football-api/internal/domains/teams/routes"
)

func main() {
	cfg := config.Load()
	db := database.Init(cfg)

	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.Recovery())
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.CORS(cfg))

	api := r.Group("/api/v1")
	{
		routes.Register(api, db, cfg)
		teamroutes.Register(api, db, cfg)
		playerroutes.Register(api, db, cfg)
		matchroutes.Register(api, db, cfg)
		reportroutes.Register(api, db, cfg)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": true, "message": "ok"})
	})

	log.Printf("server running on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
