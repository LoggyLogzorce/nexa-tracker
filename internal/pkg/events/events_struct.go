package events

import (
	"github.com/google/uuid"
	"time"
)

// ProjectEvent данные события удаления проекта
type ProjectEvent struct {
	Type      EventType
	ProjectID uuid.UUID
}

// ToEvent конвертирует ProjectEvent в events.Event
func (e ProjectEvent) ToEvent() Event {
	return Event{
		Type:      e.Type,
		Timestamp: time.Now(),
		Data:      e,
	}
}

// TaskEvent данные события для задачи
type TaskEvent struct {
	Type      EventType
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Title       string     `gorm:"not null;size:100" json:"title"`
	Description *string    `gorm:"type:text" json:"description"`
	Deadline    *time.Time `json:"deadline"`

	ProjectID  uuid.UUID  `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE" json:"project_id"`
	StatusID   *uint      `json:"status_id"`
	PriorityID *uint      `json:"priority_id"`
	AssigneeID *uuid.UUID `gorm:"type:uuid" json:"assignee_id"`
	ReporterID *uuid.UUID `gorm:"type:uuid" json:"reporter_id"`
	DeletedBy  uuid.UUID  `gorm:"type:uuid" json:"deleted_by"`
}

// ToEvent конвертирует ProjectEvent в events.Event
func (e TaskEvent) ToEvent() Event {
	return Event{
		Type:      e.Type,
		Timestamp: time.Now(),
		Data:      e,
	}
}
