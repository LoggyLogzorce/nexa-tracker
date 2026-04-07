package participant

import "github.com/google/uuid"

type ProjectParticipant struct {
	ProjectID uint      `gorm:"primaryKey" json:"project_id"`
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	Role      string    `gorm:"not null;size:10" json:"role"` // owner, member, read_only
}
