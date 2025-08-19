package services

import (
	"github.com/surajNirala/go_cliemr/internal/repository"
	"github.com/surajNirala/go_cliemr/pkg/utils"
)

type AuthService struct {
	repo *repository.AuthRepository
}

func NewAuthService(repo *repository.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Login(req *utils.LoginRequest) (string, string, error) {
	return s.repo.Login(req)
}

func (s *AuthService) Refresh(req *utils.RefreshTokenRequest) (string, error) {
	return s.repo.Refresh(req)
}

func (s *AuthService) Logout(req *utils.RefreshTokenRequest) (string, error) {
	return s.repo.Logout(req)
}
