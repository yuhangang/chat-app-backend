package tables

import (
	"time"
)

type User struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	Username  string     `gorm:"type:varchar(100);not null;index" json:"username"`
	ChatRooms []ChatRoom `gorm:"foreignKey:UserID" json:"chat_rooms"`
}

type ChatRoom struct {
	ID           uint          `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	Name         string        `gorm:"type:varchar(100)" json:"name"`
	UserID       uint          `gorm:"not null;index" json:"user_id"`
	ChatMessages []ChatMessage `gorm:"foreignKey:ChatRoomID" json:"chat_messages"`
}

type ChatMessage struct {
	ID             uint             `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time        `gorm:"autoCreateTime" json:"created_at"`
	Body           string           `gorm:"type:text;not null" json:"body"`
	ChatRoomID     uint             `gorm:"not null;index" json:"chat_room_id"`
	IsUser         bool             `gorm:"not null" json:"is_user"`
	HasAttachments bool             `gorm:"default:false" json:"has_attachments"`
	Attachments    []ChatAttachment `gorm:"foreignKey:MessageID" json:"attachments"`
}

type ChatAttachment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	FileName  string    `gorm:"type:varchar(255);not null" json:"file_name"`
	FileType  string    `gorm:"type:varchar(50);not null" json:"file_type"`  // e.g., image/png, application/pdf
	FileSize  int64     `gorm:"not null" json:"file_size"`                   // File size in bytes
	FilePath  string    `gorm:"type:varchar(255);not null" json:"file_path"` // Path or URL to the file
	MessageID uint      `gorm:"not null;index" json:"message_id"`            // Foreign key to ChatMessage
}

type LlmModel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	ModelKey  string    `gorm:"type:varchar(100);not null" json:"model_key"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Creator   string    `gorm:"type:varchar(100);not null" json:"creator"`
}
