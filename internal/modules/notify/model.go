package notify

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	UserID           uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	Type             string     `gorm:"not null;size:50" json:"type"` // task_assigned, task_updated, comment_added, etc.
	Title            string     `gorm:"not null;size:100" json:"title"`
	Message          string     `gorm:"not null;type:text" json:"message"`
	RelatedTaskID    *uint      `json:"related_task_id,omitempty"`
	RelatedProjectID *uuid.UUID `gorm:"type:uuid" json:"related_project_id,omitempty"`
	IsRead           bool       `gorm:"default:false" json:"is_read"`
}
