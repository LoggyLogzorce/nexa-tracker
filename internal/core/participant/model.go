package participant

import (
	"github.com/google/uuid"
)

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
