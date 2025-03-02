package repository

import (
	"context"

	"github.com/yuhangang/chat-app-backend/internal/db/tables"
	"gorm.io/gorm"
)

type ChatConfigRepositoryImpl struct {
	conn *gorm.DB
}

func NewChatConfigRepo(conn *gorm.DB) *ChatConfigRepositoryImpl {
	return &ChatConfigRepositoryImpl{conn: conn}
}

func (repo *ChatConfigRepositoryImpl) GetChatModels(ctx context.Context) ([]tables.LlmModel, error) {
	var chatModels []tables.LlmModel

	err := repo.conn.Find(&chatModels).Error

	return chatModels, err
}
