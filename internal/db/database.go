package db

import (
	"os"

	"github.com/yuhangang/chat-app-backend/internal/db/tables"
	"github.com/yuhangang/chat-app-backend/internal/log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	var db *gorm.DB
	dbFile := "database.db"
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
	//db.Migrator().DropTable(&tables.User{}, &tables.ChatRoom{}, &tables.ChatMessage{}, &tables.ChatAttachment{})

	// Ensure the table exists before running queries
	err = db.AutoMigrate(&tables.User{}, &tables.ChatRoom{}, &tables.ChatMessage{}, &tables.ChatAttachment{}, &tables.LlmModel{})

	if err != nil {
		log.ErrorLogger.Fatalf("Failed to migrate database: %v", err)
	}

	/// seed llm models
	models := []tables.LlmModel{
		{
			ModelKey:  "gemini-2.0-flash",
			Name:      "Gemini 2.0 Flash",
			Creator:   "Google",
			Available: true,
		},
		{
			ModelKey:  "gemini-1.5-flash",
			Name:      "Gemini 1.5 Flash",
			Creator:   "OpenAI",
			Available: true,
		},
		{
			ModelKey:  "gemini-2.0-flash-thinking-exp-01-21",
			Name:      "Gemini 2.0 Flash Thinking Exp 01-21",
			Creator:   "OpenAI",
			Available: false,
		},
		{
			ModelKey:  "gemini-2.0-pro-exp-02-05",
			Name:      "Gemini 2.0 Pro Exp 02-05",
			Creator:   "OpenAI",
			Available: false,
		},
	}

	err = db.Create(&models).Error

	if err != nil {
		log.ErrorLogger.Fatalf("Failed to seed database with llm models: %v", err)
	}

	return db, nil
}
