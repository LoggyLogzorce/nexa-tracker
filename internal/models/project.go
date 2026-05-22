package models

import (
	"github.com/google/uuid"
	"time"
)

type Project struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	Title       string    `gorm:"not null;size:50" json:"title"`
	Description *string   `gorm:"size:255" json:"description"`
	OwnerID     uuid.UUID `gorm:"type:uuid;not null" json:"owner_id"`
	Status      *string   `gorm:"size:20;default:'plan'" json:"status"`
	Priority    *string   `gorm:"size:10;default:'medium'" json:"priority"`
	Owner       User      `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE" json:"-"`
}
