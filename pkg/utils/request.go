package utils

import (
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func CheckPasswordHash(req_password, user_password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user_password), []byte(req_password))
	return err == nil
}
