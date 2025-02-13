package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/yuhangang/chat-app-backend/internal/db"
	"github.com/yuhangang/chat-app-backend/internal/db/tables"
	"github.com/yuhangang/chat-app-backend/pkg/ctxkey"
)

type ChatHandlerImpl struct {
	chatRepository db.ChatRepository
}

func NewChatHandler(
	chatRepository db.ChatRepository,
) *ChatHandlerImpl {
	return &ChatHandlerImpl{
		chatRepository: chatRepository,
	}
}

func (h *ChatHandlerImpl) GetChatRoom(w http.ResponseWriter, r *http.Request) {

	// get chat room ID from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "chatID is required", http.StatusBadRequest)
	}

	chatRoomIDStr := parts[2]

	if chatRoomIDStr == "" {
		http.Error(w, "chatID is required", http.StatusBadRequest)
		return
	}

	chatRoomID, err := strconv.ParseUint(chatRoomIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid chatID", http.StatusBadRequest)
		return
	}

	chatRoom, err := h.chatRepository.GetRoomByID(r.Context(), uint(chatRoomID))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	userID := r.Context().Value(ctxkey.UserIDKey).(uint)
	hasAcessToChatRoom := h.userHasAccessToChatRoom(userID, chatRoom)

	if !hasAcessToChatRoom {
		http.Error(w, "user does not have access to chat room", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	response, err := json.Marshal(chatRoom)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (h *ChatHandlerImpl) GetChatRooms(w http.ResponseWriter, r *http.Request) {
	// get all chat rooms for the user
	userID := r.Context().Value(ctxkey.UserIDKey).(uint)
	chatRooms, err := h.chatRepository.GetChatRoomsForUser(r.Context(), userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	response, err := json.Marshal(chatRooms)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (h *ChatHandlerImpl) userHasAccessToChatRoom(userID uint, chatRoom tables.ChatRoom) bool {
	return chatRoom.UserID == userID
}

func (h *ChatHandlerImpl) DeleteChatRoom(w http.ResponseWriter, r *http.Request) {

	// get chat room ID from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "chatID is required", http.StatusBadRequest)
	}

	chatRoomIDStr := parts[2]

	if chatRoomIDStr == "" {
		http.Error(w, "chatID is required", http.StatusBadRequest)
		return
	}

	chatRoomID, err := strconv.ParseUint(chatRoomIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid chatID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(ctxkey.UserIDKey).(uint)
	err = h.chatRepository.DeleteRoomByID(r.Context(), uint(chatRoomID), userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
