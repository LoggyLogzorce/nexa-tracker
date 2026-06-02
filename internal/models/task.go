package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"time"
)

type Task struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Title       string     `gorm:"not null;size:100" json:"title"`
	Description *string    `gorm:"type:text" json:"description"`
	Deadline    *time.Time `json:"deadline"`

	ProjectID  uuid.UUID  `gorm:"type:uuid;not null" json:"project_id"`
	StatusID   *uint      `json:"status_id"`
	PriorityID *uint      `json:"priority_id"`
	AssigneeID *uuid.UUID `gorm:"type:uuid" json:"assignee_id"`
	ReporterID *uuid.UUID `gorm:"type:uuid" json:"reporter_id"`
	Project    Project    `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"-"`
	Status     *Status    `gorm:"foreignKey:StatusID;constraint:OnDelete:SET NULL" json:"-"`
	Priority   *Priority  `gorm:"foreignKey:PriorityID;constraint:OnDelete:SET NULL" json:"-"`
	Assignee   *User      `gorm:"foreignKey:AssigneeID;constraint:OnDelete:SET NULL" json:"-"`
	Reporter   *User      `gorm:"foreignKey:ReporterID;constraint:OnDelete:SET NULL" json:"-"`

	IsArchive bool `gorm:"type:bool" json:"is_archive"`
}

type UpdateHistory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	UserID  *uuid.UUID     `gorm:"type:uuid" json:"user_id"`
	TaskID  uint           `gorm:"not null" json:"task_id"`
	Old     datatypes.JSON `gorm:"type:jsonb" json:"old"`
	New     datatypes.JSON `gorm:"type:jsonb" json:"new"`
	Changes datatypes.JSON `gorm:"type:jsonb" json:"changes"`
	User    User           `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"-"`
	Task    Task           `gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE" json:"-"`
}
