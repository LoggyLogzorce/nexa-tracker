package task

import (
	"gorm.io/datatypes"
	"nexa-task-tracker/internal/core/priority"
	"nexa-task-tracker/internal/core/project"
	"nexa-task-tracker/internal/core/status"
	"nexa-task-tracker/internal/core/user"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Title       string     `gorm:"not null;size:100" json:"title"`
	Description *string    `gorm:"type:text" json:"description"`
	Deadline    *time.Time `json:"deadline"`

	ProjectID  uuid.UUID          `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE" json:"project_id"`
	StatusID   *uint              `json:"status_id"`
	PriorityID *uint              `json:"priority_id"`
	AssigneeID *uuid.UUID         `gorm:"type:uuid" json:"assignee_id"`
	ReporterID *uuid.UUID         `gorm:"type:uuid" json:"reporter_id"`
	Project    project.Project    `gorm:"foreignKey:ProjectID" json:"-"`
	Status     *status.Status     `gorm:"foreignKey:StatusID;constraint:OnDelete:SET NULL" json:"-"`
	Priority   *priority.Priority `gorm:"foreignKey:PriorityID;constraint:OnDelete:SET NULL" json:"-"`
	Assignee   *user.User         `gorm:"foreignKey:AssigneeID;constraint:OnDelete:SET NULL" json:"-"`
	Reporter   *user.User         `gorm:"foreignKey:ReporterID;constraint:OnDelete:SET NULL" json:"-"`

	IsArchive bool `gorm:"type:bool" json:"is_archive"`
}

type TaskResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Title       string  `json:"title"`
	Description *string `json:"description"`
	Deadline    *string `json:"deadline"`

	ProjectID uuid.UUID `json:"project_id"`

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

type UpdateHistory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	UserID  uuid.UUID      `gorm:"type:uuid;constraint:OnDelete:SET NULL" json:"user_id"`
	TaskID  uint           `gorm:"not null;constraint:OnDelete:CASCADE" json:"task_id"`
	Old     datatypes.JSON `gorm:"type:jsonb" json:"old"`
	New     datatypes.JSON `gorm:"type:jsonb" json:"new"`
	Changes datatypes.JSON `gorm:"type:jsonb" json:"changes"`
	User    user.User      `gorm:"foreignKey:UserID" json:"-"`
	Task    Task           `gorm:"foreignKey:TaskID" json:"-"`
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
