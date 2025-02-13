package repository

import (
	"context"

	"github.com/yuhangang/chat-app-backend/internal/db/tables"

	"gorm.io/gorm"
)

type ChatRoomRepo struct {
	conn *gorm.DB
}

func NewChatRoomRepo(conn *gorm.DB) *ChatRoomRepo {
	return &ChatRoomRepo{conn: conn}
}

func (repo *ChatRoomRepo) GetChatRoomsForUser(ctx context.Context, userID uint) ([]tables.ChatRoom, error) {
	var chatRooms []tables.ChatRoom

	err := repo.conn.WithContext(ctx).Where("user_id = ?", userID).Find(&chatRooms).Error

	return chatRooms, err
}

func (repo *ChatRoomRepo) DeleteRoomByID(ctx context.Context, chatRoomID uint, userID uint) error {
	err := repo.conn.WithContext(ctx).Where("id = ? AND user_id = ?", chatRoomID, userID).Delete(&tables.ChatRoom{}).Error

	return err
}

func (repo *ChatRoomRepo) CheckChatRoomExists(ctx context.Context, chatRoomID uint) (bool, error) {
	var chatRoom tables.ChatRoom

	err := repo.conn.WithContext(ctx).Where("id = ?", chatRoomID).First(&chatRoom).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (repo *ChatRoomRepo) GetRoomByID(ctx context.Context, chatRoomID uint) (tables.ChatRoom, error) {
	var chatRoom tables.ChatRoom

	repo.conn.Preload("ChatMessages.Attachments").First(&chatRoom, chatRoomID)
	err := repo.conn.WithContext(ctx).Where("id = ?", chatRoomID).First(&chatRoom).Error

	return chatRoom, err
}
