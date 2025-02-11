package db

import (
	"example/user/hello/internal/db/tables"
	"example/user/hello/internal/log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	var db *gorm.DB
	dbFile := "notifications.db"
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		file, err := os.Create(dbFile)
		if err != nil {
			log.ErrorLogger.Fatalf("Failed to create database file: %v", err)
		}
		file.Close()
	}

	var err error
	db, err = gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		log.ErrorLogger.Fatalf("Failed to connect to the database: %v", err)
	}

	// reset the database
	//	db.Migrator().DropTable(&tables.ChatRoom{}, &tables.ChatMessage{})

	// Ensure the table exists before running queries
	err = db.AutoMigrate(&tables.User{}, &tables.ChatRoom{}, &tables.ChatMessage{}, &tables.ChatAttachment{})

	if err != nil {
		log.ErrorLogger.Fatalf("Failed to migrate database: %v", err)
	}

	return db, nil
}
