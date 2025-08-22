package container

import (
	"github.com/surajNirala/go_cliemr/internal/handlers"
	"github.com/surajNirala/go_cliemr/internal/repository"
	"github.com/surajNirala/go_cliemr/internal/services"
	"gorm.io/gorm"
)

type Container struct {
	UserHandler  *handlers.UserHandler
	AuthHandler  *handlers.AuthHandler
	ExcelHandler *handlers.ExcelHandler

	// add more handlers here
}

func NewContainer(db *gorm.DB) *Container {
	// Repositories
	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(db)
	excelRepo := repository.NewExcelRepository(db)

	// Services
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(authRepo)
	excelService := services.NewExcelService(excelRepo)

	// Handlers
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)
	excelHandler := handlers.NewExcelHandler(excelService, db)

	return &Container{
		UserHandler:  userHandler,
		AuthHandler:  authHandler,
		ExcelHandler: excelHandler,
	}
}
