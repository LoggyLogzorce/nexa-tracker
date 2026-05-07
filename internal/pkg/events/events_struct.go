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
