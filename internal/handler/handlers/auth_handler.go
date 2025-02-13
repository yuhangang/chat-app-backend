package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/yuhangang/chat-app-backend/internal/db"
	"github.com/yuhangang/chat-app-backend/internal/db/tables"
	"github.com/yuhangang/chat-app-backend/internal/handler"
	"github.com/yuhangang/chat-app-backend/types"
)

type AuthHandlerImpl struct {
	userRepository db.UserRepository
	jwtService     types.JwtService
}

func NewAuthHandler(userRepo db.UserRepository, jwtService types.JwtService) *AuthHandlerImpl {
	return &AuthHandlerImpl{
		userRepository: userRepo,
		jwtService:     jwtService,
	}
}

// RefreshToken is a handler that refreshes the access token
func (h *AuthHandlerImpl) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Extract refresh token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "missing authorization header", http.StatusUnauthorized)
		return
	}

	// Check if it starts with "Bearer "
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
		return
	}

	// Extract the token
	refreshToken := authHeader[7:]
	if refreshToken == "" {
		http.Error(w, "missing refresh token", http.StatusUnauthorized)
		return
	}

	// Get new access token
	accessToken, err := h.jwtService.RefreshAccessToken(refreshToken)

	if err != nil {
		switch err {
		case handler.ErrExpiredToken:
			http.Error(w, "refresh token expired", http.StatusUnauthorized)
		case handler.ErrInvalidToken:
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Create response
	response := map[string]string{
		"access_token": accessToken,
	}

	// Convert to JSON and write response
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "failed to create response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (h *AuthHandlerImpl) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Get username or generate one if empty
	username := r.URL.Query().Get("username")
	if username == "" {
		username = "user_" + uuid.New().String()[:8] // Generates a short unique username
	}

	user := tables.User{
		Username: username,
	}

	userCreated, err := h.userRepository.CreateUser(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate JWT
	jwtPayload, err := h.jwtService.GenerateTokens(userCreated.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return JWT and user as JSON
	jwtPayloadResponse := UserResponse{
		AccessToken:  jwtPayload.AccessToken,
		RefreshToken: jwtPayload.RefreshToken,
		User:         userCreated,
	}

	w.Header().Set("Content-Type", "application/json") // Correct header order
	w.WriteHeader(http.StatusCreated)

	// Return user as JSON
	if err := json.NewEncoder(w).Encode(jwtPayloadResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
