package services

import (
	"strings"

	"golang.org/x/crypto/bcrypt"

	"football-api/internal/domains/users/dtos"
	usermodels "football-api/internal/domains/users/models"
	"football-api/internal/domains/users/repositories"
)

type Service interface {
	Create(req dtos.CreateUserRequest) (*dtos.UserResponse, error)
	List(offset, limit int, search string) ([]dtos.UserResponse, int64, error)
}

type service struct {
	repo repositories.Repository
}

func New(repo repositories.Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(req dtos.CreateUserRequest) (*dtos.UserResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := usermodels.User{
		Name:         strings.TrimSpace(req.Name),
		Email:        strings.ToLower(strings.TrimSpace(req.Email)),
		PasswordHash: string(hash),
		Role:         req.Role,
	}
	if err := s.repo.Create(&user); err != nil {
		return nil, err
	}
	return &dtos.UserResponse{ID: user.ID, Name: user.Name, Email: user.Email, Role: user.Role}, nil
}

func (s *service) List(offset, limit int, search string) ([]dtos.UserResponse, int64, error) {
	users, total, err := s.repo.List(offset, limit, search)
	if err != nil {
		return nil, 0, err
	}
	out := make([]dtos.UserResponse, 0, len(users))
	for _, u := range users {
		out = append(out, dtos.UserResponse{ID: u.ID, Name: u.Name, Email: u.Email, Role: u.Role})
	}
	return out, total, nil
}
