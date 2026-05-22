package task

import (
	"gorm.io/datatypes"
	"time"

	"github.com/google/uuid"
)

type TaskResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Title       string  `json:"title"`
	Description *string `json:"description"`
	Deadline    *string `json:"deadline"`

	ProjectID    uuid.UUID `json:"project_id"`
	ProjectTitle string    `json:"project_title"`

	Status   *TaskStatusResponse   `json:"status"`
	Priority *TaskPriorityResponse `json:"priority"`

	Assignee *TaskUserResponse `json:"assignee"`
	Reporter *TaskUserResponse `json:"reporter"`

	IsArchive bool `json:"is_archive"`
}

type TaskStatusResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	OrderIndex int    `json:"order_index"`
}

type TaskPriorityResponse struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Color string `json:"color"`
}

type TaskUserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	AvatarUrl string    `json:"avatar_url"`
}

type HistoryResponse struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	User    TaskUserResponse `gorm:"type:uuid" json:"user"`
	TaskID  uint             `gorm:"not null" json:"task_id"`
	Old     datatypes.JSON   `gorm:"type:jsonb" json:"old"`
	New     datatypes.JSON   `gorm:"type:jsonb" json:"new"`
	Changes datatypes.JSON   `gorm:"type:jsonb" json:"changes"` // [{field, old_value, new_value}]
}
