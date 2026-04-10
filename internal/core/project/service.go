package project

import (
	"context"
	"github.com/google/uuid"
	"nexa-task-tracker/internal/core/priority"
	"nexa-task-tracker/internal/core/status"
	"time"
)

type Service interface {
	Create(ctx context.Context, project *Project, ownerID uuid.UUID) error
	GetByID(ctx context.Context, id uint) (*Project, error)
	List(ctx context.Context, userID uuid.UUID) ([]Project, error)
	Update(ctx context.Context, project *Project) error
	Delete(ctx context.Context, id uint) error
}

type service struct {
	repo         Repository
	statusRepo   status.Repository
	priorityRepo priority.Repository
}

func NewService(repo Repository, statusRepo status.Repository, priorityRepo priority.Repository) Service {
	return &service{
		repo:         repo,
		statusRepo:   statusRepo,
		priorityRepo: priorityRepo,
	}
}

func (s *service) Create(ctx context.Context, project *Project, ownerID uuid.UUID) error {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 1. Установить ownerID
	project.OwnerID = ownerID

	// 2. Создать проект
	if err := s.repo.Create(ctxT, project); err != nil {
		return err
	}

	// 3. Создать дефолтные статусы
	defaultStatuses := []status.Status{
		{ProjectID: project.ID, Name: "To Do", Color: "#808080", OrderIndex: 0},
		{ProjectID: project.ID, Name: "In Progress", Color: "#3b82f6", OrderIndex: 1},
		{ProjectID: project.ID, Name: "Done", Color: "#22c55e", OrderIndex: 2},
	}

	defaultPriorities := []priority.Priority{
		{ProjectID: project.ID, Title: "Low", Color: "#808080"},
		{ProjectID: project.ID, Title: "Medium", Color: "#3b82f6"},
		{ProjectID: project.ID, Title: "High", Color: "#22c55e"},
	}

	if err := s.statusRepo.CreateBatch(defaultStatuses); err != nil {
		return err
	}

	if err := s.priorityRepo.CreateBatch(defaultPriorities); err != nil {
		return err
	}

	return nil
}

func (s *service) GetByID(ctx context.Context, id uint) (*Project, error) {
	// TODO: Implement
	return nil, nil
}

func (s *service) List(ctx context.Context, userID uuid.UUID) ([]Project, error) {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	projects, err := s.repo.List(ctxT, userID)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (s *service) Update(ctx context.Context, project *Project) error {
	// TODO: Implement
	return nil
}

func (s *service) Delete(ctx context.Context, id uint) error {
	// TODO: Implement
	return nil
}
