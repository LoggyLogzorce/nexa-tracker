package status

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"nexa-task-tracker/internal/core/project"
	"nexa-task-tracker/internal/pkg/events"
	"time"
)

type Service interface {
	Create(ctx context.Context, status *Status) error
	GetByID(ctx context.Context, id uint) (*Status, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]Status, error)
	Update(ctx context.Context, id uint, projectID uuid.UUID, updates UpdateStatusRequest) (*Status, error)
	Delete(ctx context.Context, id uint, projectID uuid.UUID) error
	HandleProjectDeleted(event events.Event) error
	HandleProjectCreated(event events.Event) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, status *Status) error {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 1. Проверка уникальности названия
	_, err := s.repo.GetByName(ctxT, status.Name)
	if err == nil {
		return ErrStatusNameExists
	}

	// 2. Валидация color (если указан)
	if status.Color != "" {
		if err := ValidateHexColor(status.Color); err != nil {
			return err
		}
	} else {
		// Установить дефолтное значение
		status.Color = "#cccccc"
	}

	// 3. Если order_index не указан (равен 0), установить в конец
	if status.OrderIndex == 0 {
		maxIndex, err := s.repo.GetMaxOrderIndex(ctxT, status.ProjectID)
		if err != nil {
			return err
		}
		status.OrderIndex = maxIndex + 1
	} else {
		// Проверить, что статус с таким order_index не существует
		existingStatuses, err := s.repo.GetByProjectID(ctxT, status.ProjectID)
		if err != nil {
			return err
		}
		for _, existing := range existingStatuses {
			if existing.OrderIndex == status.OrderIndex {
				return ErrDuplicateOrderIndex
			}
		}
	}

	// 4. Создать статус
	if err := s.repo.Create(ctxT, status); err != nil {
		return err
	}

	return nil
}

func (s *service) GetByID(ctx context.Context, id uint) (*Status, error) {
	// TODO: Implement
	return nil, nil
}

func (s *service) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]Status, error) {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 1. Получить статусы проекта
	statuses, err := s.repo.GetByProjectID(ctxT, projectID)
	if err != nil {
		return nil, err
	}

	return statuses, nil
}

func (s *service) Update(ctx context.Context, id uint, projectID uuid.UUID, updates UpdateStatusRequest) (*Status, error) {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 1. Получить существующий статус
	status, err := s.repo.GetByID(ctxT, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStatusNotFound
		}
		return nil, err
	}

	// 2. Проверить, что статус принадлежит указанному проекту
	if status.ProjectID != projectID {
		return nil, ErrStatusNotFound
	}

	// 3. Обработать обновление name (если указано)
	if updates.Name != nil {
		// Проверить уникальность имени (исключая текущий статус)
		statuses, err := s.repo.GetByProjectID(ctxT, projectID)
		if err != nil {
			return nil, err
		}
		for _, existing := range statuses {
			if existing.Name == *updates.Name && existing.ID != id {
				return nil, ErrStatusNameExists
			}
		}
		status.Name = *updates.Name
	}

	// 4. Обработать обновление color (если указано)
	if updates.Color != nil {
		if err := ValidateHexColor(*updates.Color); err != nil {
			return nil, err
		}
		status.Color = *updates.Color
	}

	// 5. Обработать обновление order_index (если указано)
	if updates.OrderIndex != nil {
		newOrderIndex := *updates.OrderIndex
		oldOrderIndex := status.OrderIndex

		// Получить все статусы проекта
		statuses, err := s.repo.GetByProjectID(ctxT, projectID)
		if err != nil {
			return nil, err
		}

		// Проверить максимальный order_index
		maxOrderIndex := len(statuses) - 1
		if newOrderIndex > maxOrderIndex {
			// Если превышает, поставить последним
			newOrderIndex = maxOrderIndex
		}

		if newOrderIndex != oldOrderIndex {
			// Подготовить массовое обновление для сдвига других статусов
			var batchUpdates []struct {
				ID         uint
				OrderIndex int
			}

			if newOrderIndex < oldOrderIndex {
				// Перемещение вверх: сдвинуть вниз статусы между new и old
				for _, st := range statuses {
					if st.ID != id && st.OrderIndex >= newOrderIndex && st.OrderIndex < oldOrderIndex {
						batchUpdates = append(batchUpdates, struct {
							ID         uint
							OrderIndex int
						}{
							ID:         st.ID,
							OrderIndex: st.OrderIndex + 1,
						})
					}
				}
			} else {
				// Перемещение вниз: сдвинуть вверх статусы между old и new
				for _, st := range statuses {
					if st.ID != id && st.OrderIndex > oldOrderIndex && st.OrderIndex <= newOrderIndex {
						batchUpdates = append(batchUpdates, struct {
							ID         uint
							OrderIndex int
						}{
							ID:         st.ID,
							OrderIndex: st.OrderIndex - 1,
						})
					}
				}
			}

			// Выполнить массовое обновление
			if len(batchUpdates) > 0 {
				if err := s.repo.UpdateOrderIndexBatch(ctxT, batchUpdates); err != nil {
					return nil, err
				}
			}

			status.OrderIndex = newOrderIndex
		}
	}

	// 6. Сохранить обновленный статус
	if err := s.repo.Update(ctxT, status); err != nil {
		return nil, err
	}

	return status, nil
}

func (s *service) Delete(ctx context.Context, id uint, projectID uuid.UUID) error {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 1. Получить статус для проверки существования и принадлежности проекту
	status, err := s.repo.GetByID(ctxT, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrStatusNotFound
		}
		return err
	}

	// 2. Проверить, что статус принадлежит указанному проекту
	if status.ProjectID != projectID {
		return ErrStatusNotFound
	}

	deletedOrderIndex := status.OrderIndex

	// 3. Удалить статус
	if err := s.repo.Delete(ctxT, id); err != nil {
		return err
	}

	// 4. Получить все оставшиеся статусы проекта
	statuses, err := s.repo.GetByProjectID(ctxT, projectID)
	if err != nil {
		return err
	}

	// 5. Пересчитать order_index для статусов, которые были после удаленного
	var batchUpdates []struct {
		ID         uint
		OrderIndex int
	}

	for _, st := range statuses {
		if st.OrderIndex > deletedOrderIndex {
			batchUpdates = append(batchUpdates, struct {
				ID         uint
				OrderIndex int
			}{
				ID:         st.ID,
				OrderIndex: st.OrderIndex - 1,
			})
		}
	}

	// 6. Выполнить массовое обновление order_index
	if len(batchUpdates) > 0 {
		if err := s.repo.UpdateOrderIndexBatch(ctxT, batchUpdates); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) HandleProjectDeleted(event events.Event) error {
	data, ok := event.Data.(project.ProjectEvent)
	if !ok {
		return fmt.Errorf("invalid event data type")
	}

	ctxT, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Удалить все статусы проекта
	if err := s.repo.DeleteByProjectID(ctxT, data.ProjectID); err != nil {
		return fmt.Errorf("failed to delete statuses for project %s: %w", data.ProjectID, err)
	}

	return nil
}

func (s *service) HandleProjectCreated(event events.Event) error {
	data, ok := event.Data.(project.ProjectEvent)
	if !ok {
		return fmt.Errorf("invalid event data type")
	}

	defaultStatuses := []Status{
		{ProjectID: data.ProjectID, Name: "To Do", Color: "#808080", OrderIndex: 0},
		{ProjectID: data.ProjectID, Name: "In Progress", Color: "#3b82f6", OrderIndex: 1},
		{ProjectID: data.ProjectID, Name: "Done", Color: "#22c55e", OrderIndex: 2},
	}

	ctxT, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.repo.CreateBatch(ctxT, defaultStatuses); err != nil {
		return fmt.Errorf("failed to create default statuses for project %s: %w", data.ProjectID, err)
	}

	return nil
}
