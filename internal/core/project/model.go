package project

import (
	"nexa-task-tracker/internal/models"
	"time"

	"github.com/google/uuid"
)

type ProjectResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`

	CreatedAt time.Time `json:"created_at"`
	Owner     struct {
		ID        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		AvatarUrl string    `json:"avatar_url"`
	} `json:"owner"`
	Status     *string           `json:"status"`
	Priority   *string           `json:"priority"`
	UserRole   string            `json:"user_role"`
	Statuses   []models.Status   `json:"statuses"`
	Priorities []models.Priority `json:"priorities"`
}
