package status

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"nexa-task-tracker/internal/core/project"
	"nexa-task-tracker/internal/pkg/events"
	"time"
)

type Service interface {
	Create(ctx context.Context, status *Status) error
	GetByID(ctx context.Context, id uint) (*Status, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]Status, error)
	Update(ctx context.Context, status *Status) error
	Delete(ctx context.Context, id uint) error
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

func (s *service) Update(ctx context.Context, status *Status) error {
	// TODO: Implement
	return nil
}

func (s *service) Delete(ctx context.Context, id uint) error {
	// TODO: Implement
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
