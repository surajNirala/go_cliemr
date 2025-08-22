package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/surajNirala/go_cliemr/internal/container"
	"github.com/surajNirala/go_cliemr/internal/middleware"
)

func InitRoutes(r *gin.Engine, c *container.Container) {

	// Now attach routes
	v1 := r.Group("/api/v1") // versioning is common practice
	{
		v1.POST("/login", c.AuthHandler.Login)
		v1.POST("/refresh", c.AuthHandler.Refresh)
		v1.POST("/logout", c.AuthHandler.Logout)
		v1.GET("/users", middleware.AuthMiddleware(), middleware.RequireRoles(1, 2), c.UserHandler.GetUserList)
		v1.POST("/patient-import", middleware.AuthMiddleware(), middleware.RequireRoles(1, 2), c.ExcelHandler.PatientImport)
	}
}
