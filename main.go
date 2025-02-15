package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// env of gemini api key
	log.Println("Starting server...")

	httpServer := NewHttpServer(ctx, ":8080")
	fileServer := NewFileServer(":3002")
	go func() {
		if err := httpServer.Run(); err != nil {
			log.Fatalf("HTTP server failed to start: %v", err)
		}
	}()

	go func() {
		if err := fileServer.Run(); err != nil {
			log.Fatalf("File server failed to start: %v", err)
		}
	}()

	select {}

}
