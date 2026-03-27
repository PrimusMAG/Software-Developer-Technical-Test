package database

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"football-api/internal/core/config"
	matchmodels "football-api/internal/domains/matches/models"
	playermodels "football-api/internal/domains/players/models"
	teammodels "football-api/internal/domains/teams/models"
	usermodels "football-api/internal/domains/users/models"
	"football-api/internal/shared/constants"
)

func Init(cfg config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("failed connect db: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed get sql db: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := db.AutoMigrate(
		&usermodels.User{},
		&teammodels.Team{},
		&playermodels.Player{},
		&matchmodels.Match{},
		&matchmodels.GoalEvent{},
	); err != nil {
		log.Fatalf("failed migrate db: %v", err)
	}

	seed(db)
	return db
}

func seed(db *gorm.DB) {
	var count int64
	db.Model(&usermodels.User{}).Where("email = ?", "admin@xyz.com").Count(&count)
	if count == 0 {
		hash, _ := bcrypt.GenerateFromPassword([]byte("Admin123!"), bcrypt.DefaultCost)
		_ = db.Create(&usermodels.User{
			Name:         "Admin XYZ",
			Email:        "admin@xyz.com",
			PasswordHash: string(hash),
			Role:         string(constants.RoleAdmin),
		}).Error
	}

	db.Model(&teammodels.Team{}).Count(&count)
	if count == 0 {
		teams := []teammodels.Team{
			{Name: "Garuda Muda", FoundedYear: 2012, HQAddress: "Jl. Merdeka No.1", HQCity: "Jakarta"},
			{Name: "Rajawali FC", FoundedYear: 2015, HQAddress: "Jl. Pemuda No.7", HQCity: "Bandung"},
		}
		_ = db.Create(&teams).Error

		var t1, t2 teammodels.Team
		_ = db.Where("name = ?", "Garuda Muda").First(&t1).Error
		_ = db.Where("name = ?", "Rajawali FC").First(&t2).Error

		players := []playermodels.Player{
			{TeamID: t1.ID, Name: "Budi", HeightCM: 178, WeightKG: 72, Position: string(constants.PositionForward), JerseyNumber: 9},
			{TeamID: t1.ID, Name: "Andi", HeightCM: 180, WeightKG: 75, Position: string(constants.PositionDefender), JerseyNumber: 4},
			{TeamID: t2.ID, Name: "Rizky", HeightCM: 176, WeightKG: 70, Position: string(constants.PositionMidfielder), JerseyNumber: 8},
			{TeamID: t2.ID, Name: "Dimas", HeightCM: 182, WeightKG: 78, Position: string(constants.PositionGoalkeeper), JerseyNumber: 1},
		}
		_ = db.Create(&players).Error
	}
}
