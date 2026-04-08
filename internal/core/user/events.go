package user

import (
	"time"

	"github.com/google/uuid"
	"nexa-task-tracker/internal/pkg/events"
)

// UserDeletedEvent данные события удаления пользователя
type UserDeletedEvent struct {
	UserID uuid.UUID
	Email  string
	Name   string
}

// ToEvent конвертирует UserDeletedEvent в events.Event
func (e UserDeletedEvent) ToEvent() events.Event {
	return events.Event{
		Type:      events.UserDeleted,
		Timestamp: time.Now(),
		Data:      e,
	}
}
