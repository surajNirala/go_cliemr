package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/surajNirala/go_cliemr/internal/services"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUserList(c *gin.Context) {

	users, err := h.service.GetUserList()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Users not found",
			"issue": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Fetched User list",
		"data":    users,
	})
}
