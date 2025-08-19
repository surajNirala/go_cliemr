package config

import (
	"fmt"
	"log"
	"os"

	"github.com/surajNirala/go_cliemr/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	// DB = db
	fmt.Println("✅ Connected to MySQL")
	// Run migrations
	err = db.AutoMigrate(&models.RefreshToken{})
	if err != nil {
		log.Fatal("failed migration:", err)
	}

	log.Println("Migration completed ✅")
	return db, nil
}
