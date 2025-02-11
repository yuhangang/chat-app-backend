package repository

import (
	"context"
	"example/user/hello/internal/db/tables"
	"example/user/hello/internal/service"
	"mime/multipart"

	"gorm.io/gorm"
)

type MessageRepo struct {
	conn           *gorm.DB
	storageService service.StorageService
}

func NewMessageRepo(conn *gorm.DB, storageService service.StorageService) *MessageRepo {
	return &MessageRepo{conn: conn, storageService: storageService}
}

func (repo *MessageRepo) GetMessagesForChatRoom(ctx context.Context, chatRoomID uint) ([]tables.ChatMessage, error) {
	var chatMessages []tables.ChatMessage

	err := repo.conn.WithContext(ctx).Where("chat_room_id = ?", chatRoomID).Find(&chatMessages).Error

	return chatMessages, err
}

func (repo *MessageRepo) CreateMessage(
	ctx context.Context,
	chatRoomID uint,
	message string,
	response string,
	attachment *multipart.FileHeader) ([]tables.ChatMessage, error) {
	// Create the chat message for the user
	chatMessage := tables.ChatMessage{
		ChatRoomID:     chatRoomID,
		Body:           message,
		IsUser:         true,
		HasAttachments: attachment != nil,
	}

	// Create the chat message for the response
	chatResponse := tables.ChatMessage{
		ChatRoomID: chatRoomID,
		Body:       response,
	}

	// Start a transaction to ensure both message and attachments are saved atomically
	err := repo.conn.Transaction(func(tx *gorm.DB) error {
		var err error

		// Save the user message
		err = tx.WithContext(ctx).Create(&chatMessage).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		// Save the bot response message
		err = tx.WithContext(ctx).Create(&chatResponse).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		if attachment != nil {
			file, err := attachment.Open()
			if err != nil {
				tx.Rollback()
				return err
			}
			defer file.Close()

			// Save the file to disk or cloud storage
			filePath, err := repo.storageService.SaveFile(attachment)
			if err != nil {
				tx.Rollback()
				return err
			}

			// Create the attachment record
			attachment := tables.ChatAttachment{
				FileName:  attachment.Filename,
				FileType:  attachment.Header.Get("Content-Type"),
				FileSize:  attachment.Size,
				FilePath:  "http://localhost:3002/" + filePath,
				MessageID: chatMessage.ID, // Attach to the user's message
			}

			// Save the attachment to the database
			err = tx.WithContext(ctx).Create(&attachment).Error
			if err != nil {
				tx.Rollback()
				return err
			}

			chatMessage.Attachments = append(chatMessage.Attachments, attachment)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Return the user message and the response
	return []tables.ChatMessage{chatMessage, chatResponse}, nil
}

func (repo *MessageRepo) CreateChatRoomWithMessage(
	ctx context.Context,
	userID uint,
	chatRoomName string,
	message string,
	response string,
	attachment *multipart.FileHeader) (tables.ChatRoom, error) {
	// Create the chat room
	chatRoom := tables.ChatRoom{
		UserID: userID,
		Name:   chatRoomName,
	}

	// Start a transaction to ensure both message and attachments are saved atomically
	err := repo.conn.Transaction(func(tx *gorm.DB) error {
		var err error

		err = repo.conn.WithContext(ctx).Create(&chatRoom).Error

		if err != nil {
			return err
		}

		chatMessage := tables.ChatMessage{
			ChatRoomID:     chatRoom.ID,
			Body:           message,
			IsUser:         true,
			HasAttachments: attachment != nil,
		}

		// Create the chat message for the response
		chatResponse := tables.ChatMessage{
			ChatRoomID: chatRoom.ID,
			Body:       response,
		}

		// Save the user message
		err = tx.WithContext(ctx).Create(&chatMessage).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		// Save the bot response message
		err = tx.WithContext(ctx).Create(&chatResponse).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		if attachment != nil {
			file, err := attachment.Open()
			if err != nil {
				tx.Rollback()
				return err
			}
			defer file.Close()

			// Save the file to disk or cloud storage
			filePath, err := repo.storageService.SaveFile(attachment)
			if err != nil {
				tx.Rollback()
				return err
			}

			// Create the attachment record
			attachment := tables.ChatAttachment{
				FileName:  attachment.Filename,
				FileType:  attachment.Header.Get("Content-Type"),
				FileSize:  attachment.Size,
				FilePath:  "http://localhost:3002/" + filePath,
				MessageID: chatMessage.ID, // Attach to the user's message
			}

			// Save the attachment to the database
			err = tx.WithContext(ctx).Create(&attachment).Error
			if err != nil {
				tx.Rollback()
				return err
			}

			chatMessage.Attachments = append(chatMessage.Attachments, attachment)
		}

		chatRoom.ChatMessages = append(chatRoom.ChatMessages, chatMessage, chatResponse)

		return nil
	})

	return chatRoom, err
}
