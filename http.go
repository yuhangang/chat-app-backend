package main

import (
	"context"
	"log"
	"net/http"

	"github.com/yuhangang/chat-app-backend/internal/db"
	"github.com/yuhangang/chat-app-backend/internal/db/repository"
	"github.com/yuhangang/chat-app-backend/internal/handler"
	"github.com/yuhangang/chat-app-backend/internal/handler/handlers"
	"github.com/yuhangang/chat-app-backend/internal/service/services/gemini_service"
	"github.com/yuhangang/chat-app-backend/internal/service/services/jwt_service"
	"github.com/yuhangang/chat-app-backend/internal/service/services/storage_service"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type httpServer struct {
	addr        string
	httpHandler *handler.Handler
}

func NewHttpServer(ctx context.Context, addr string) *httpServer {

	conn, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		panic(err)
	}

	llmService := gemini_service.NewGeminiServiceV1(ctx)
	jwtService, err := jwt_service.NewJwtService()
	storageService := storage_service.NewStorageServiceV1()

	if err != nil {
		log.Fatalf("Failed to create jwt service: %v", err)
		panic(err)
	}

	userRepository := repository.NewUserRepo(conn)
	chatRepository := repository.NewChatRoomRepo(conn)
	chatConfigRepository := repository.NewChatConfigRepo(conn)
	messageRepo := repository.NewMessageRepo(conn, storageService)
	llmRepo := repository.NewLLMRepo(conn, llmService)

	userHandler := handlers.NewUserHandler(userRepository)
	chatHandler := handlers.NewChatHandler(chatRepository)
	chatConfigHandler := handlers.NewChatConfigHandler(chatConfigRepository)
	messageHandler := handlers.NewMessageChatHandler(chatRepository, messageRepo, llmRepo)
	authHandler := handlers.NewAuthHandler(userRepository, jwtService)

	httpHandler := handler.NewHandler(chatHandler, chatConfigHandler, messageHandler, userHandler, authHandler, jwtService)

	return &httpServer{addr: addr, httpHandler: httpHandler}
}

func (s *httpServer) Run() error {
	router := mux.NewRouter()
	s.httpHandler.RegisterRoutes(router)

	log.Println("Starting server on", s.addr)

	// CORS Middleware should be applied before starting the server
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://example.com", "http://localhost:3000"}, // Allow specific domains
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	// Start the server with CORS middleware applied

	return http.ListenAndServe(":8080", c.Handler(router))
}
