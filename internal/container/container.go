package container

import (
	"github.com/surajNirala/go_cliemr/internal/handlers"
	"github.com/surajNirala/go_cliemr/internal/repository"
	"github.com/surajNirala/go_cliemr/internal/services"
	"gorm.io/gorm"
)

type Container struct {
	UserHandler *handlers.UserHandler
	// add more handlers here
}

func NewContainer(db *gorm.DB) *Container {
	// Repositories
	userRepo := repository.NewUserRepository(db)

	// Services
	userService := services.NewUserService(userRepo)

	// Handlers
	userHandler := handlers.NewUserHandler(userService)

	return &Container{
		UserHandler: userHandler,
	}
}
