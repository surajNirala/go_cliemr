package container

import (
	"github.com/surajNirala/go_cliemr/internal/handlers"
	"github.com/surajNirala/go_cliemr/internal/repository"
	"github.com/surajNirala/go_cliemr/internal/services"
	"gorm.io/gorm"
)

type Container struct {
	UserHandler *handlers.UserHandler
	AuthHandler *handlers.AuthHandler
	// add more handlers here
}

func NewContainer(db *gorm.DB) *Container {
	// Repositories
	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(db)

	// Services
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(authRepo)

	// Handlers
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)

	return &Container{
		UserHandler: userHandler,
		AuthHandler: authHandler,
	}
}
