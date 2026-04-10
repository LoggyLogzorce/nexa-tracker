package priority

import (
	"fmt"
	"github.com/google/uuid"
	"nexa-task-tracker/internal/core/project"
	"nexa-task-tracker/internal/pkg/events"
)

type Service interface {
	Create(priority *Priority) error
	GetByID(id uint) (*Priority, error)
	GetByProjectID(projectID uuid.UUID) ([]Priority, error)
	Update(priority *Priority) error
	Delete(id uint) error
	HandleProjectCreated(event events.Event) error
	HandleProjectDeleted(event events.Event) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(priority *Priority) error {
	// TODO: Implement
	return nil
}

func (s *service) GetByID(id uint) (*Priority, error) {
	// TODO: Implement
	return nil, nil
}

func (s *service) GetByProjectID(projectID uuid.UUID) ([]Priority, error) {
	// TODO: Implement
	return nil, nil
}

func (s *service) Update(priority *Priority) error {
	// TODO: Implement
	return nil
}

func (s *service) Delete(id uint) error {
	// TODO: Implement
	return nil
}

func (s *service) HandleProjectDeleted(event events.Event) error {
	data, ok := event.Data.(project.ProjectEvent)
	if !ok {
		return fmt.Errorf("invalid event data type")
	}

	// Удалить все статусы проекта
	if err := s.repo.DeleteByProjectID(data.ProjectID); err != nil {
		return fmt.Errorf("failed to delete statuses for project %s: %w", data.ProjectID, err)
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

	if err := s.repo.CreateBatch(defaultPriorities); err != nil {
		return fmt.Errorf("failed to create default statuses for project %s: %w", data.ProjectID, err)
	}

	return nil
}
