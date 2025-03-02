package repository

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/google/uuid"
	"github.com/yuhangang/chat-app-backend/internal/db/tables"
	"github.com/yuhangang/chat-app-backend/types"
	"gorm.io/gorm"
)

type LLMRepo struct {
	conn       *gorm.DB
	llmService types.HttpServiceV1
}

func NewLLMRepo(conn *gorm.DB, llmService types.HttpServiceV1) *LLMRepo {
	return &LLMRepo{conn: conn, llmService: llmService}
}

func (r *LLMRepo) CallGemini(ctx context.Context, prompt string, chatroomId uint, file *multipart.FileHeader,
) (types.GeminiApiResponse, error) {
	var history []*genai.Content
	var err error

	var sessionId string
	if chatroomId == 0 {
		sessionId = generateSessionID()

	} else {
		sessionId = r.getSessionId(chatroomId)
		history, err = r.getChatHistory(chatroomId)
		if err != nil {
			return types.GeminiApiResponse{}, err
		}
	}

	if file != nil {
		tempFilePath, err := saveUploadedFile(file)

		if err != nil {
			return types.GeminiApiResponse{}, err
		}
		defer os.Remove(tempFilePath)

		return r.llmService.SendFileWithText(ctx, sessionId, prompt, history, tempFilePath)
	}

	return r.llmService.CallGemini(ctx, sessionId, prompt, history)
}

func saveUploadedFile(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	tempFile, err := os.CreateTemp("", fileHeader.Filename)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy file to temp file: %w", err)
	}

	return tempFile.Name(), nil
}

func (r *LLMRepo) getChatHistory(chatroomId uint) ([]*genai.Content, error) {
	var messages []struct {
		Body   string
		IsUser bool
	}
	err := r.conn.Model(&tables.ChatMessage{}).
		Where("chat_room_id = ?", chatroomId).
		Scan(&messages).Error

	history := make([]*genai.Content, len(messages))
	for i, m := range messages {
		history[i] = &genai.Content{
			Parts: []genai.Part{
				genai.Text(m.Body),
			},
			Role: func() string {
				if m.IsUser {
					return "user"
				}
				return "model"
			}(),
		}
	}

	return history, err
}

func generateSessionID() string {
	return uuid.New().String()
}

func (r *LLMRepo) getSessionId(chatroomId uint) string {
	var chatRoom tables.ChatRoom
	r.conn.Model(&tables.ChatRoom{}).Select("session_id").First(&chatRoom, chatroomId)

	return chatRoom.SessionID

}
