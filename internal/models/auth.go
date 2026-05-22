package models

import (
	"github.com/google/uuid"
	"time"
)

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	UserID    uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	TokenHash string     `gorm:"uniqueIndex;not null;size:255" json:"-"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	UserAgent *string    `gorm:"size:255" json:"user_agent,omitempty"`
	IPAddress *string    `gorm:"type:inet" json:"ip_address,omitempty"`
	User      User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}
