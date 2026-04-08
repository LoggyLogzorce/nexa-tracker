package notify

import (
	"fmt"
	"nexa-task-tracker/internal/core/user"
	"nexa-task-tracker/internal/pkg/events"
)

// HandleUserDeleted обрабатывает событие удаления пользователя
func (s *service) HandleUserDeleted(event events.Event) error {
	// Извлечь данные события
	data, ok := event.Data.(user.UserDeletedEvent)
	if !ok {
		return fmt.Errorf("invalid event data type")
	}

	// Удалить все уведомления пользователя
	if err := s.repo.DeleteByUserID(data.UserID); err != nil {
		return fmt.Errorf("failed to delete notifications for user %s: %w", data.UserID, err)
	}

	return nil
}
