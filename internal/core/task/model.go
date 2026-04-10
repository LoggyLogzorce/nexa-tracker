package task

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Title       string     `gorm:"not null;size:100" json:"title"`
	Description *string    `gorm:"type:text" json:"description"`
	Deadline    *time.Time `json:"deadline"`

	ProjectID  uuid.UUID  `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE" json:"project_id"`
	StatusID   *uint      `json:"status_id"`
	PriorityID *uint      `json:"priority_id"`
	AssigneeID *uuid.UUID `gorm:"type:uuid" json:"assignee_id"`
	ReporterID *uuid.UUID `gorm:"type:uuid" json:"reporter_id"`
}

// BeforeUpdate hook for logging changes to update_history
func (t *Task) BeforeUpdate(tx *gorm.DB) error {
	// TODO: Implement update history logging
	// 1. Get old task state from DB
	// 2. Compare with new state
	// 3. Create UpdateHistory record with old/new JSONB
	return nil
}
