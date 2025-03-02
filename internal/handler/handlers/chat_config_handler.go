package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/yuhangang/chat-app-backend/internal/db"
)

type ChatConfigHandlerImpl struct {
	chatConfigRepo db.ChatConfigRepository
}

func NewChatConfigHandler(chatConfigRepo db.ChatConfigRepository) *ChatConfigHandlerImpl {
	return &ChatConfigHandlerImpl{chatConfigRepo: chatConfigRepo}
}

func (h *ChatConfigHandlerImpl) GetChatModels(w http.ResponseWriter, r *http.Request) {
	chatModels, err := h.chatConfigRepo.GetChatModels(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chatModels)
}
