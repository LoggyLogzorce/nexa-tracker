package project

import (
	"nexa-task-tracker/internal/core/priority"
	"nexa-task-tracker/internal/core/status"
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	Title       string    `gorm:"not null;size:50" json:"title"`
	Description *string   `gorm:"size:255" json:"description"`
	OwnerID     uuid.UUID `gorm:"type:uuid;not null" json:"owner_id"`
}

type ProjectResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`

	CreatedAt time.Time `json:"created_at"`
	Owner     struct {
		ID    uuid.UUID `json:"id"`
		Name  string    `json:"name"`
		Email string    `json:"email"`
	} `json:"owner"`
	Statuses   []status.Status     `json:"statuses"`
	Priorities []priority.Priority `json:"priorities"`
}
