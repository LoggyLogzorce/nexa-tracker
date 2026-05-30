package models

import (
	"github.com/google/uuid"
	"time"
)

type UserProvider struct {
	ID uint `gorm:"primaryKey"`

	UserID uuid.UUID `gorm:"type:uuid;not null;index"`

	Provider       string `gorm:"size:50;not null"`
	ProviderUserID string `gorm:"size:255;not null"`

	CreatedAt time.Time

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (UserProvider) TableName() string {
	return "user_providers"
}
