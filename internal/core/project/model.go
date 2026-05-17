package project

import (
	"nexa-task-tracker/internal/core/priority"
	"nexa-task-tracker/internal/core/status"
	"nexa-task-tracker/internal/core/user"
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	Title       string    `gorm:"not null;size:50" json:"title"`
	Description *string   `gorm:"size:255" json:"description"`
	OwnerID     uuid.UUID `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE" json:"owner_id"`
	Status      *string   `gorm:"size:20;default:'plan'" json:"status"`
	Priority    *string   `gorm:"size:10;default:'medium'" json:"priority"`
	Owner       user.User `gorm:"foreignKey:OwnerID" json:"-"`
}

type ProjectResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`

	CreatedAt time.Time `json:"created_at"`
	Owner     struct {
		ID        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		AvatarUrl string    `json:"avatar_url"`
	} `json:"owner"`
	Status     *string             `json:"status"`
	Priority   *string             `json:"priority"`
	UserRole   string              `json:"user_role"`
	Statuses   []status.Status     `json:"statuses"`
	Priorities []priority.Priority `json:"priorities"`
}
