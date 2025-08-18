package services

import (
	"github.com/surajNirala/go_cliemr/internal/models"
	"github.com/surajNirala/go_cliemr/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserList() ([]*models.User, error) {
	return s.repo.GetUserList()
}
