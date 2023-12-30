package bootstrap

import (
	"fmt"
	"go_sample_api/internal/course"
	"go_sample_api/internal/user"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connection(logger *log.Logger) (*gorm.DB, error) {
	var db *gorm.DB
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Printf("Database not connected - %s", err.Error())
		return nil, err
	}
	logger.Println("Database connected successfully")

	if os.Getenv("DATABASE_DEBUG") == "true" {
		db = db.Debug()
		logger.Println("Database run in mode debug")
	}

	if os.Getenv("DATABASE_MIGRATE") == "true" {
		if err := db.AutoMigrate(&user.User{}); err != nil {
			return nil, err
		}

		if err := db.AutoMigrate(&course.Course{}); err != nil {
			return nil, err
		}
		logger.Println("Database migrated")
	}

	return db, nil
}

func InitLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
}
