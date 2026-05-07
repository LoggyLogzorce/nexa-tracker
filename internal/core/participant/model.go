package participant

import "github.com/google/uuid"

type ProjectParticipant struct {
	ProjectID uuid.UUID `gorm:"type:uuid;primaryKey;constraint:OnDelete:CASCADE" json:"project_id"`
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	Role      string    `gorm:"not null;size:10" json:"role"` // owner, member, read_only
}

type ProjectParticipantsResponse struct {
	ProjectID uuid.UUID `json:"project_id"`
	Role      string    `json:"role"`
	User      struct {
		UserID uuid.UUID `json:"user_id"`
		Name   string    `json:"name"`
		Email  string    `json:"email"`
	}
}
