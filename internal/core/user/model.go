package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	Email        string  `gorm:"uniqueIndex;not null;size:255" json:"email"`
	PasswordHash string  `gorm:"not null;size:255" json:"-"`
	Name         string  `gorm:"not null;size:50" json:"name"`
	Role         string  `gorm:"not null;default:'user';size:20" json:"role"` // admin, user
	Secret2FA    *string `gorm:"size:255" json:"-"`
}
