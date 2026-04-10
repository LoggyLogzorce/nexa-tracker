package status

import (
	"fmt"
	"github.com/google/uuid"
	"nexa-task-tracker/internal/core/project"
	"nexa-task-tracker/internal/pkg/events"
)

type Service interface {
	Create(status *Status) error
	GetByID(id uint) (*Status, error)
	GetByProjectID(projectID uuid.UUID) ([]Status, error)
	Update(status *Status) error
	Delete(id uint) error
	HandleProjectDeleted(event events.Event) error
	HandleProjectCreated(event events.Event) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(status *Status) error {
	// TODO: Implement
	return nil
}

func (s *service) GetByID(id uint) (*Status, error) {
	// TODO: Implement
	return nil, nil
}

func (s *service) GetByProjectID(projectID uuid.UUID) ([]Status, error) {
	// TODO: Implement
	return nil, nil
}

func (s *service) Update(status *Status) error {
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

	defaultStatuses := []Status{
		{ProjectID: data.ProjectID, Name: "To Do", Color: "#808080", OrderIndex: 0},
		{ProjectID: data.ProjectID, Name: "In Progress", Color: "#3b82f6", OrderIndex: 1},
		{ProjectID: data.ProjectID, Name: "Done", Color: "#22c55e", OrderIndex: 2},
	}

	if err := s.repo.CreateBatch(defaultStatuses); err != nil {
		return fmt.Errorf("failed to create default statuses for project %s: %w", data.ProjectID, err)
	}

	return nil
}
