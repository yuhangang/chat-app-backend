package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/yuhangang/chat-app-backend/internal/db"
	"github.com/yuhangang/chat-app-backend/pkg/ctxkey"
)

type MessageHandlerImpl struct {
	messageRepository db.MessageRepository
	chatRepository    db.ChatRepository
	llmRepository     db.LLMRepository
}

func NewMessageChatHandler(
	chatRepository db.ChatRepository,
	messageRepository db.MessageRepository,
	llmRepository db.LLMRepository,
) *MessageHandlerImpl {
	return &MessageHandlerImpl{
		chatRepository:    chatRepository,
		messageRepository: messageRepository,
		llmRepository:     llmRepository,
	}
}

func (h *MessageHandlerImpl) CreateChatRoomWithMessage(w http.ResponseWriter, r *http.Request) {
	var err error
	// read the request ['prompt'] from the request
	prompt := r.FormValue("prompt")
	_, fileHeader, _ := r.FormFile("attachment")

	userId := r.Context().Value(ctxkey.UserIDKey).(uint)

	geminiResponse, err := h.llmRepository.CallGemini(r.Context(), prompt, 0, false, fileHeader)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// use first 5 words of the prompt as the chat room name
	words := strings.Split(prompt, " ")
	var chatRoomName string
	if len(words) < 10 {
		chatRoomName = strings.Join(words, " ")
	} else {
		chatRoomName = strings.Join(words[:10], " ")
		chatRoomName = chatRoomName + "..."
	}
	chatRoom, err := h.messageRepository.CreateChatRoomWithMessage(r.Context(), userId, chatRoomName, prompt, geminiResponse.Response, fileHeader)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	response, err := json.Marshal(chatRoom)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (h *MessageHandlerImpl) CreateMessage(w http.ResponseWriter, r *http.Request) {
	prompt := r.FormValue("prompt")

	// Get chat room ID from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "chatID is required", http.StatusBadRequest)
		return
	}
	chatRoomIDStr := parts[2]
	chatRoomID, err := strconv.ParseUint(chatRoomIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid chatID", http.StatusBadRequest)
		return
	}

	// Check if chat room exists
	exists, err := h.chatRepository.CheckChatRoomExists(r.Context(), uint(chatRoomID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "chat room does not exist", http.StatusBadRequest)
		return
	}

	// Get the uploaded files (attachments)
	_, fileHeader, _ := r.FormFile("attachment")

	// Call Gemini for a response based on the prompt
	geminiResponse, err := h.llmRepository.CallGemini(r.Context(), prompt, uint(chatRoomID), true, fileHeader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create message and attachments in the repository
	createdMessage, err := h.messageRepository.CreateMessage(r.Context(), uint(chatRoomID), prompt, geminiResponse.Response, fileHeader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response, err := json.Marshal(createdMessage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}
