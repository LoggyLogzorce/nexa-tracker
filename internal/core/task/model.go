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

	IsArchive bool `gorm:"type:bool" json:"is_archive"`
}

type TaskResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Title       string  `json:"title"`
	Description *string `json:"description"`
	Deadline    *string `json:"deadline"`

	ProjectID uuid.UUID `json:"project_id"`

	Status   *TaskStatusResponse   `json:"status"`
	Priority *TaskPriorityResponse `json:"priority"`

	Assignee *TaskUserResponse `json:"assignee"`
	Reporter *TaskUserResponse `json:"reporter"`

	IsArchive bool `json:"is_archive"`
}

// BeforeUpdate hook for logging changes to update_history
func (t *Task) BeforeUpdate(tx *gorm.DB) error {
	// TODO: Implement update history logging
	// 1. Get old task state from DB
	// 2. Compare with new state
	// 3. Create UpdateHistory record with old/new JSONB
	return nil
}

type TaskStatusResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	OrderIndex int    `json:"order_index"`
}

type TaskPriorityResponse struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Color string `json:"color"`
}

type TaskUserResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}
