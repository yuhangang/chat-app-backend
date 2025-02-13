package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/yuhangang/chat-app-backend/pkg/ctxkey"
	"github.com/yuhangang/chat-app-backend/types"

	"github.com/gorilla/mux"
)

type Handler struct {
	chatHandler    ChatHandler
	messageHandler MessageHandler
	userHandler    UserHandler
	authHandler    AuthHandler
	jwtService     types.JwtService
}

func NewHandler(chatHandler ChatHandler, messageHandler MessageHandler, userHandler UserHandler, authHandler AuthHandler, jwtService types.JwtService) *Handler {
	return &Handler{
		chatHandler:    chatHandler,
		messageHandler: messageHandler,
		userHandler:    userHandler,
		authHandler:    authHandler,
		jwtService:     jwtService,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	// Only JWT required
	jwtProtectedRoutes := map[string]func(http.ResponseWriter, *http.Request){
		"POST /chats":        h.messageHandler.CreateChatRoomWithMessage,
		"GET /chats":         h.chatHandler.GetChatRooms,
		"GET /chats/{id}":    h.chatHandler.GetChatRoom,
		"DELETE /chats/{id}": h.chatHandler.DeleteChatRoom,
		"POST /chats/{id}":   h.messageHandler.CreateMessage,
		"GET /user":          h.userHandler.GetUser,
	}

	// No protection
	publicRoutes := map[string]func(http.ResponseWriter, *http.Request){
		"POST /auth":         h.authHandler.CreateUser,
		"POST /auth/refresh": h.authHandler.RefreshToken,
	}

	for route, handler := range jwtProtectedRoutes {
		parts := strings.Split(route, " ")
		method, path := parts[0], parts[1]
		router.HandleFunc(path, h.jwtAuthMiddleware(handler)).Methods(method)
	}

	for route, handler := range publicRoutes {
		parts := strings.Split(route, " ")
		method, path := parts[0], parts[1]
		router.HandleFunc(path, handler).Methods(method)
	}

}

/*
func (h *Handler) fullAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serverKey := r.Header.Get("X-Server-Key")
		if serverKey != h.jwtService.serverKey {
			http.Error(w, "Invalid server key", http.StatusUnauthorized)
			return
		}

		h.jwtAuthMiddleware(next)(w, r)
	}
}
*/

func (h *Handler) jwtAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		claims, err := h.jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ctxkey.UserIDKey, claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

type ChatHandler interface {
	GetChatRoom(http.ResponseWriter, *http.Request)
	GetChatRooms(http.ResponseWriter, *http.Request)
	DeleteChatRoom(http.ResponseWriter, *http.Request)
}

type UserHandler interface {
	GetUser(http.ResponseWriter, *http.Request)
}

type MessageHandler interface {
	CreateChatRoomWithMessage(http.ResponseWriter, *http.Request)
	CreateMessage(http.ResponseWriter, *http.Request)
}

type AuthHandler interface {
	CreateUser(http.ResponseWriter, *http.Request)
	RefreshToken(http.ResponseWriter, *http.Request)
}

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("token has expired")
	ErrMissingJWTSecret = errors.New("missing JWT secrets in environment")
)
