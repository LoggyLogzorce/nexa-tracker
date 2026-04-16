package priority

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
	Create(ctx context.Context, priority *Priority) error
	GetByID(ctx context.Context, id uint) (*Priority, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]Priority, error)
	Update(ctx context.Context, id uint, projectID uuid.UUID, updates UpdatePriorityRequest) (*Priority, error)
	Delete(ctx context.Context, id uint, projectID uuid.UUID) error
	HandleProjectCreated(event events.Event) error
	HandleProjectDeleted(event events.Event) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, priority *Priority) error {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 1. Проверка уникальности title в рамках проекта
	_, err := s.repo.GetByTitle(ctxT, priority.Title, priority.ProjectID)
	if err == nil {
		return ErrPriorityTitleExists
	}

	// 2. Валидация color (если указан)
	if priority.Color != "" {
		if err := ValidateHexColor(priority.Color); err != nil {
			return err
		}
	} else {
		// Установить дефолтное значение
		priority.Color = "#cccccc"
	}

	// 3. Создать приоритет
	if err := s.repo.Create(ctxT, priority); err != nil {
		return err
	}

	return nil
}

func (s *service) GetByID(ctx context.Context, id uint) (*Priority, error) {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	priority, err := s.repo.GetByID(ctxT, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPriorityNotFound
		}
		return nil, err
	}

	return priority, nil
}

func (s *service) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]Priority, error) {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Получить приоритеты проекта
	priorities, err := s.repo.GetByProjectID(ctxT, projectID)
	if err != nil {
		return nil, err
	}

	return priorities, nil
}

func (s *service) Update(ctx context.Context, id uint, projectID uuid.UUID, updates UpdatePriorityRequest) (*Priority, error) {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 1. Получить существующий приоритет
	priority, err := s.repo.GetByID(ctxT, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPriorityNotFound
		}
		return nil, err
	}

	// 2. Проверить, что приоритет принадлежит указанному проекту
	if priority.ProjectID != projectID {
		return nil, ErrPriorityNotFound
	}

	// 3. Обработать обновление title (если указано)
	if updates.Title != nil {
		// Проверить уникальность title (исключая текущий приоритет)
		existingPriority, err := s.repo.GetByTitle(ctxT, *updates.Title, projectID)
		if err == nil && existingPriority.ID != id {
			return nil, ErrPriorityTitleExists
		}
		priority.Title = *updates.Title
	}

	// 4. Обработать обновление color (если указано)
	if updates.Color != nil {
		if err := ValidateHexColor(*updates.Color); err != nil {
			return nil, err
		}
		priority.Color = *updates.Color
	}

	// 5. Сохранить обновленный приоритет
	if err := s.repo.Update(ctxT, priority); err != nil {
		return nil, err
	}

	return priority, nil
}

func (s *service) Delete(ctx context.Context, id uint, projectID uuid.UUID) error {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 1. Получить приоритет для проверки существования и принадлежности проекту
	priority, err := s.repo.GetByID(ctxT, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPriorityNotFound
		}
		return err
	}

	// 2. Проверить, что приоритет принадлежит указанному проекту
	if priority.ProjectID != projectID {
		return ErrPriorityNotFound
	}

	// 3. Удалить приоритет (задачи автоматически получат NULL благодаря constraint)
	if err := s.repo.Delete(ctxT, id); err != nil {
		return err
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

	// Удалить все приоритеты проекта
	if err := s.repo.DeleteByProjectID(ctxT, data.ProjectID); err != nil {
		return fmt.Errorf("failed to delete priorities for project %s: %w", data.ProjectID, err)
	}

	return nil
}

func (s *service) HandleProjectCreated(event events.Event) error {
	data, ok := event.Data.(project.ProjectEvent)
	if !ok {
		return fmt.Errorf("invalid event data type")
	}

	defaultPriorities := []Priority{
		{ProjectID: data.ProjectID, Title: "Low", Color: "#22c55e"},    // зелёный
		{ProjectID: data.ProjectID, Title: "Medium", Color: "#f59e0b"}, // жёлтый/оранжевый
		{ProjectID: data.ProjectID, Title: "High", Color: "#ef4444"},   // красный
	}

	ctxT, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.repo.CreateBatch(ctxT, defaultPriorities); err != nil {
		return fmt.Errorf("failed to create default priorities for project %s: %w", data.ProjectID, err)
	}

	return nil
}
