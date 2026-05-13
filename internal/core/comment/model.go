package comment

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	UserID  uuid.UUID `gorm:"type:uuid" json:"user_id"`
	TaskID  uint      `gorm:"not null" json:"task_id"`
	Content string    `gorm:"not null;type:text" json:"content"`
}

type CommentResponse struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	User    CommentUserResponse `gorm:"type:uuid" json:"user"`
	TaskID  uint                `gorm:"not null" json:"task_id"`
	Content string              `gorm:"not null;type:text" json:"content"`
}

type CommentUserResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}
