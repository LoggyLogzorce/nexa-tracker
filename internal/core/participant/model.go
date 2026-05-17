package participant

import (
	"nexa-task-tracker/internal/core/user"

	"github.com/google/uuid"
)

type ProjectParticipant struct {
	ProjectID uuid.UUID `gorm:"type:uuid;primaryKey;constraint:OnDelete:CASCADE" json:"project_id"`
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey;constraint:OnDelete:CASCADE" json:"user_id"`
	Role      string    `gorm:"not null;size:10" json:"role"` // owner, member, read_only
	User      user.User `gorm:"foreignKey:UserID" json:"-"`
}

type ProjectParticipantsResponse struct {
	ProjectID uuid.UUID `json:"project_id"`
	Role      string    `json:"role"`
	User      struct {
		UserID    uuid.UUID `json:"user_id"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		AvatarUrl string    `json:"avatar_url"`
	}
}
