package db

import (
	"log"
	"os"

	"github.com/Ckefa/ckefablog.git/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	dsn := os.Getenv("DSN")

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	err = sqlDB.Ping()
	if err != nil {
		return err
	}

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Println("<<Connected to DB>>")
	return nil
}
