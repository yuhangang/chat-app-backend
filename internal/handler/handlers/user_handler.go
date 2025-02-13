package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/yuhangang/chat-app-backend/internal/db"
	"github.com/yuhangang/chat-app-backend/internal/db/tables"
	"github.com/yuhangang/chat-app-backend/pkg/ctxkey"
)

type UserHandlerImpl struct {
	userRepository db.UserRepository
}

func NewUserHandler(userRepo db.UserRepository) *UserHandlerImpl {
	return &UserHandlerImpl{
		userRepository: userRepo,
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
