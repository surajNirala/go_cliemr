package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/surajNirala/go_cliemr/internal/services"
	"github.com/surajNirala/go_cliemr/pkg/utils"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// ---------------LOGIN---------------
func (h *AuthHandler) Login(c *gin.Context) {
	var req utils.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			utils.RespondError(c, "error", http.StatusBadRequest, "All Fields are required.", utils.ValidationError(errs))
			return
		}
		utils.RespondError(c, "error", http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}
	// utils.RequestAll(c)
	accessToken, refreshToken, err := h.service.Login(&req)
	if err != nil {
		utils.RespondError(c, "error", http.StatusNotFound, err.Error(), err.Error())
		return
	}
	data := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
	utils.RespondSuccess(c, "success", http.StatusOK, "Login successful", data)
}

// -------- REFRESH --------
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req utils.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
		message := "Refresh token required"
		utils.RespondError(c, "error", http.StatusBadRequest, message, message)
		return
	}

	// Check DB for valid refresh token
	accessToken, err := h.service.Refresh(&req)
	if err != nil {
		utils.RespondError(c, "error", http.StatusInternalServerError, err.Error(), err.Error())
		return
	}

	data := map[string]string{
		"access_token": accessToken,
	}
	utils.RespondSuccess(c, "success", http.StatusOK, "New Access token.", data)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Implement logout logic here, if needed
	var req utils.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
		message := "Authorization token is required"
		utils.RespondError(c, "error", http.StatusBadRequest, message, message)
		return
	}
	message, err := h.service.Logout(&req)
	if err != nil {
		utils.RespondError(c, "error", http.StatusBadRequest, err.Error(), err.Error())
		return
	}
	// For now, we can just return a success message
	utils.RespondSuccess(c, "success", http.StatusOK, message, nil)
}
