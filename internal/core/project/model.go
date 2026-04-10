package project

import (
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
