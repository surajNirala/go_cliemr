package repository

import (
	"github.com/surajNirala/go_cliemr/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserList() ([]*models.User, error) {
	var users []*models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetUserListRepositoryRaw(db *gorm.DB) ([]*models.User, error) {
	var users []*models.User
	query := "SELECT id,username,email,password,created_at, updated_at FROM users"
	if err := db.Raw(query).Scan(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
