package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"football-api/internal/core/config"
	"football-api/internal/core/middleware"
	"football-api/internal/domains/auth/dtos"
	"football-api/internal/domains/auth/repositories"
)

type Service interface {
	Login(req dtos.LoginRequest) (*dtos.AuthResponse, error)
	Refresh(token string) (*dtos.AuthResponse, error)
}

type service struct {
	repo repositories.Repository
	cfg  config.Config
}

func New(repo repositories.Repository, cfg config.Config) Service {
	return &service{repo: repo, cfg: cfg}
}

func (s *service) Login(req dtos.LoginRequest) (*dtos.AuthResponse, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credential")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credential")
	}
	return s.issueTokens(user.ID, user.Role), nil
}

func (s *service) Refresh(token string) (*dtos.AuthResponse, error) {
	claims := &middleware.JWTClaims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil || !parsed.Valid {
		return nil, errors.New("invalid refresh token")
	}
	return s.issueTokens(claims.UserID, claims.Role), nil
}

func (s *service) issueTokens(userID uint, role string) *dtos.AuthResponse {
	now := time.Now()
	accessTTL := time.Duration(s.cfg.JWTAccessTTLMin) * time.Minute
	refreshTTL := time.Duration(s.cfg.JWTRefreshTTLHour) * time.Hour

	accessClaims := middleware.JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   "access",
		},
	}
	refreshClaims := middleware.JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(refreshTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   "refresh",
		},
	}

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, _ := access.SignedString([]byte(s.cfg.JWTSecret))

	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, _ := refresh.SignedString([]byte(s.cfg.JWTSecret))

	return &dtos.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresInSec: int(accessTTL.Seconds()),
		Role:         role,
	}
}
