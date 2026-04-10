package project

import (
	"github.com/google/uuid"
	"nexa-task-tracker/internal/pkg/events"
	"time"
)

// ProjectEvent данные события удаления пользователя
type ProjectEvent struct {
	Type      events.EventType
	ProjectID uuid.UUID
}

// ToEvent конвертирует ProjectEvent в events.Event
func (e ProjectEvent) ToEvent() events.Event {
	return events.Event{
		Type:      e.Type,
		Timestamp: time.Now(),
		Data:      e,
	}
}
