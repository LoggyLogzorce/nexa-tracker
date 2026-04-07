package history

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type UpdateHistory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	UserID *uuid.UUID     `gorm:"type:uuid" json:"user_id"`
	TaskID uint           `gorm:"not null" json:"task_id"`
	Old    datatypes.JSON `gorm:"type:jsonb" json:"old"`
	New    datatypes.JSON `gorm:"type:jsonb" json:"new"`
}
