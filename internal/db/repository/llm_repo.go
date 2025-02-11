package repository

import (
	"context"
	"example/user/hello/internal/db/tables"
	"example/user/hello/types"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/google/generative-ai-go/genai"
	"gorm.io/gorm"
)

type LLMRepo struct {
	conn       *gorm.DB
	llmService types.HttpServiceV1
}

func NewLLMRepo(conn *gorm.DB, llmService types.HttpServiceV1) *LLMRepo {
	return &LLMRepo{conn: conn, llmService: llmService}
}

func (r *LLMRepo) CallGemini(ctx context.Context, prompt string, chatroomId uint, useHistory bool, file *multipart.FileHeader,
) (types.GeminiApiResponse, error) {
	var history []genai.Part
	var err error
	if !useHistory {
		history, err = getChatHistory(r, chatroomId)
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

		return r.llmService.SendFileWithText(ctx, prompt, history, tempFilePath)
	}

	return r.llmService.CallGemini(ctx, prompt, history)
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

func getChatHistory(r *LLMRepo, chatroomId uint) ([]genai.Part, error) {
	var messages []string
	err := r.conn.Model(&tables.ChatMessage{}).
		Where("chat_room_id = ?", chatroomId).
		Pluck("body", &messages).Error

	history := make([]genai.Part, len(messages))
	for i, m := range messages {
		history[i] = genai.Text(m)
	}

	return history, err
}
