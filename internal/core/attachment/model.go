package attachment

import (
	"time"

	"github.com/google/uuid"
)

type Attachment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	TaskID   uint       `gorm:"not null" json:"task_id"`
	UserID   *uuid.UUID `gorm:"type:uuid" json:"user_id"`
	Filename string     `gorm:"not null;size:255" json:"filename"`
	FilePath string     `gorm:"not null;size:500" json:"file_path"`
	FileSize int64      `gorm:"not null" json:"file_size"`
	MimeType *string    `gorm:"size:100" json:"mime_type,omitempty"`
}
