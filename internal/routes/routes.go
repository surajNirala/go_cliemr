package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/surajNirala/go_cliemr/internal/container"
)

func InitRoutes(r *gin.Engine, c *container.Container) {

	// Now attach routes
	v1 := r.Group("/api/v1") // versioning is common practice
	{
		v1.GET("/users", c.UserHandler.GetUserList)
	}
}
