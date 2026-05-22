package attachment

import (
	"time"

	"github.com/google/uuid"
)

type AttachmentResponse struct {
	ID        uint                   `json:"id"`
	CreatedAt time.Time              `json:"created_at"`
	TaskID    uint                   `json:"task_id"`
	User      AttachmentUserResponse `json:"user"`
	Filename  string                 `json:"filename"`
	FileSize  int64                  `json:"file_size"`
	MimeType  *string                `json:"mime_type,omitempty"`
}

type AttachmentUserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	AvatarUrl string    `json:"avatar_url"`
}
