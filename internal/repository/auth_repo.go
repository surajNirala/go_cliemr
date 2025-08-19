package repository

import (
	"errors"
	"log"
	"time"

	"github.com/surajNirala/go_cliemr/internal/models"
	"github.com/surajNirala/go_cliemr/pkg/utils"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) Login(req *utils.LoginRequest) (string, string, error) {
	var user *models.User
	// hashPassword := utils.HashPassword(req.Password)
	result := r.db.Where("email =?", req.Email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", "", errors.New("Email/Password is incorrect")
		}
		return "", "", result.Error
	}
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return "", "", errors.New("Email/Password is incorrect")
	}
	// Generate tokens
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return "", "", errors.New("failed to generate access token")
	}
	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return "", "", errors.New("failed to generate refresh token")
	}
	// Save refresh token in DB
	refreshRecord := models.RefreshToken{
		UserID: user.ID,
		Token:  refreshToken,
		Expiry: time.Now().Add(7 * 24 * time.Hour), // 7 days expiry
	}

	if err := r.db.Create(&refreshRecord).Error; err != nil {
		return "", "", errors.New("failed to save refresh token")
	}

	return accessToken, refreshToken, nil
}

// -------- REFRESH --------
func (r *AuthRepository) Refresh(req *utils.RefreshTokenRequest) (string, error) {
	// Validate refresh token
	claims, err := utils.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	// Check DB for valid refresh token
	var stored models.RefreshToken
	if err := r.db.Where("token = ?", req.RefreshToken).First(&stored).Error; err != nil {
		return "", errors.New("Refresh token not found (maybe logged out)")
	}
	if stored.Expiry.Before(time.Now()) {
		return "", errors.New("Refresh token expired")
	}

	// Generate new access token
	accessToken, err := utils.GenerateAccessToken(claims.UserID, claims.Email)
	if err != nil {
		return "", errors.New("failed to generate new access token")
	}
	return accessToken, nil
}

// -------- LOGOUT --------
func (r *AuthRepository) Logout(req *utils.RefreshTokenRequest) (string, error) {
	// Validate refresh token
	claims, err := utils.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	// Delete refresh token from DB
	result := r.db.Where("token = ?", req.RefreshToken).Delete(&models.RefreshToken{})
	if result.Error != nil {
		return "", errors.New("failed to delete refresh token")
	}
	if result.RowsAffected == 0 {
		return "", errors.New("refresh token not found or already logged out")
	}
	log.Printf("UserID %v logged out successfully", claims.UserID)
	return "Logout successful", nil

}
