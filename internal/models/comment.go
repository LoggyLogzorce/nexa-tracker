package models

import (
	"github.com/google/uuid"
	"time"
)

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	UserID  uuid.UUID `gorm:"type:uuid" json:"user_id"`
	TaskID  uint      `gorm:"not null" json:"task_id"`
	Content string    `gorm:"not null;type:text" json:"content"`
	User    User      `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"-"`
	Task    Task      `gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE" json:"-"`
}
