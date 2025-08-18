package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/surajNirala/go_cliemr/config"
	"github.com/surajNirala/go_cliemr/internal/container"
	"github.com/surajNirala/go_cliemr/internal/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize DB
	db, err := config.InitDB()
	if err != nil {
		log.Fatal("Database error: ", err)
	}

	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Setup Gin
	r := gin.Default()

	// Build container
	c := container.NewContainer(db)
	routes.InitRoutes(r, c)
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Run server in goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %s\n", err)
		}
	}()
	log.Printf("🚀 Server running on http://localhost:%s", port)

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // wait for signal
	log.Println("Shutting down server...")

	// Context with timeout for shutdown (10 seconds max)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited gracefully")
}
