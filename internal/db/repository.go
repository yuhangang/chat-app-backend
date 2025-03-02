package db

import (
	"context"
	"mime/multipart"

	"github.com/yuhangang/chat-app-backend/internal/db/tables"
	"github.com/yuhangang/chat-app-backend/types"
)

type Database interface {
	UserRepository() UserRepository
	ChatRepository() ChatRepository
	MessageRepository() MessageRepository
}

type ChatRepository interface {
	GetChatRoomsForUser(ctx context.Context, userID uint) ([]tables.ChatRoom, error)
	GetRoomByID(ctx context.Context, chatRoomID uint) (tables.ChatRoom, error)
	DeleteRoomByID(ctx context.Context, chatRoomID uint, userID uint) error
	CheckChatRoomExists(ctx context.Context, chatRoomID uint) (bool, error)
}

type ChatConfigRepository interface {
	GetChatModels(ctx context.Context) ([]tables.LlmModel, error)
}

type MessageRepository interface {
	GetMessagesForChatRoom(ctx context.Context, chatRoomID uint) ([]tables.ChatMessage, error)
	CreateMessage(ctx context.Context,
		chatRoomID uint,
		message string,
		response string,
		attachment *multipart.FileHeader,
	) ([]tables.ChatMessage, error)
	CreateChatRoomWithMessage(
		ctx context.Context,
		userID uint,
		chatRoomName string,
		message string,
		response types.GeminiApiResponse,
		attachment *multipart.FileHeader) (tables.ChatRoom, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, user tables.User) (tables.User, error)
	GetUser(ctx context.Context, userID uint) (tables.User, error)
	GetUserByUsername(ctx context.Context, username string) (tables.User, error)
	BindUser(ctx context.Context, userID uint, username string) (tables.User, error)
}

type LLMRepository interface {
	CallGemini(ctx context.Context, prompt string, chatroomId uint, files *multipart.FileHeader,
	) (types.GeminiApiResponse, error)
}
