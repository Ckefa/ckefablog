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

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		defer mu.Unlock()
		return err
	}

	err = sqlDB.Ping()
	if err != nil {
		return err
	}

	err = DB.AutoMigrate(&models.User{}, &models.Customer{}, &models.Package{}, &models.Order{})
	if err != nil {
		log.Fatal("<<Migration Error", err)
		return err
	} else {
		log.Println("<<DB migration complete>>")
	}

	for _, p := range models.Packages {
		var existingPackage models.Package
		// Check if the package exists, if not, create it
		mu.Lock()
		result := DB.Where("id = ?", p.ID).FirstOrCreate(&existingPackage, p)
		if result.Error != nil {
			log.Fatal(result.Error)
		}
		mu.Unlock()
	}

	log.Println("<<Connected to DB>>")
	return nil
}
