package models

import (
	"github.com/google/uuid"
)

type Priority struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProjectID uuid.UUID `gorm:"type:uuid;not null" json:"project_id"`
	Title     string    `gorm:"not null;size:50" json:"title"`
	Color     string    `gorm:"default:'#cccccc';size:9" json:"color"`
	Project   Project   `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"-"`
}
