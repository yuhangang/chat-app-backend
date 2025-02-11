package handlers

import (
	"encoding/json"
	"example/user/hello/internal/db"
	"example/user/hello/internal/db/tables"
	"example/user/hello/pkg/ctxkey"
	"example/user/hello/types"
	"net/http"

	"github.com/google/uuid"
)

type UserHandlerImpl struct {
	userRepository db.UserRepository
	jwtService     types.JwtService
}

func NewUserHandler(userRepo db.UserRepository, jwtService types.JwtService) *UserHandlerImpl {
	return &UserHandlerImpl{
		userRepository: userRepo,
		jwtService:     jwtService,
	}
}

func (h *UserHandlerImpl) CreateUser(w http.ResponseWriter, r *http.Request) {
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

func (h *UserHandlerImpl) GetUser(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(ctxkey.UserIDKey).(uint)

	// Get user from DB
	user, err := h.userRepository.GetUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return user as JSON
	w.Header().Set("Content-Type", "application/json") // Correct header order
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type UserResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	User         tables.User `json:"user"`
}
