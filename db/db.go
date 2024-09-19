package db

import (
	"log"
	"os"
	"sync"

	"github.com/Ckefa/ckefablog/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var mu sync.Mutex

func Init() error {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Println("DSN not set in environment variables")
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Failed to connect to database:", err)
		return err
	}

	// Check DB connection
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	err = sqlDB.Ping()
	if err != nil {
		return err
	}

	// Auto-migrate models
	err = DB.AutoMigrate(&models.User{}, &models.Customer{}, &models.Package{}, &models.Order{})
	if err != nil {
		log.Println("Migration Error:", err)
		return err
	}
	log.Println("<<DB migration complete>>")

	// Initialize predefined packages
	for _, p := range models.Packages {
		var existingPackage models.Package

		// Lock around database insertion
		mu.Lock()
		result := DB.Where("id = ?", p.ID).FirstOrCreate(&existingPackage, p)
		mu.Unlock()

		if result.Error != nil {
			log.Println("Error inserting package:", result.Error)
			return result.Error
		}
	}

	log.Println("<<Connected to DB>>")
	return nil
}
