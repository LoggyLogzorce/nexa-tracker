package project

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"nexa-task-tracker/internal/core/priority"
	"nexa-task-tracker/internal/core/status"
	"nexa-task-tracker/internal/core/user"
	"nexa-task-tracker/internal/pkg/events"
	"time"
)

type Service interface {
	Create(ctx context.Context, project *Project, ownerID uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*ProjectResponse, error)
	List(ctx context.Context, userID uuid.UUID) ([]Project, error)
	Update(ctx context.Context, project *Project, userID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type service struct {
	repo         Repository
	eventBus     *events.EventBus
	userRepo     user.Repository
	statusRepo   status.Repository
	priorityRepo priority.Repository
}

func NewService(repo Repository, eventBus *events.EventBus, userRepo user.Repository, statusRepo status.Repository, priorityRepo priority.Repository) Service {
	return &service{
		repo:         repo,
		eventBus:     eventBus,
		userRepo:     userRepo,
		statusRepo:   statusRepo,
		priorityRepo: priorityRepo,
	}
}

func (s *service) Create(ctx context.Context, project *Project, ownerID uuid.UUID) error {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 1. Установить id и owner
	project.ID = uuid.New()
	project.OwnerID = ownerID

	// 2. Создать проект
	if err := s.repo.Create(ctxT, project); err != nil {
		return err
	}

	event := events.ProjectEvent{
		Type:      events.ProjectCreated,
		ProjectID: project.ID,
	}
	s.eventBus.Publish(event.ToEvent())

	return nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*ProjectResponse, error) {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 1. Получить проект из БД
	project, err := s.repo.GetByID(ctxT, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	projectDto := &ProjectResponse{
		ID:          project.ID,
		Title:       project.Title,
		Description: project.Description,
		CreatedAt:   project.CreatedAt,
	}

	owner, err := s.userRepo.GetByID(ctxT, project.OwnerID)
	if err != nil {
		return projectDto, ErrGetOwner
	}

	projectDto.Owner.ID = owner.ID
	projectDto.Owner.Name = owner.Name
	projectDto.Owner.Email = owner.Email

	statuses, err := s.statusRepo.GetByProjectID(ctxT, project.ID)
	if err != nil {
		return nil, err
	}

	priorities, err := s.priorityRepo.GetByProjectID(ctxT, project.ID)
	if err != nil {
		return nil, err
	}

	projectDto.Statuses = statuses
	projectDto.Priorities = priorities

	return projectDto, nil
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

func (s *service) Update(ctx context.Context, project *Project, userID uuid.UUID) error {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 1. Получить проект из БД
	existingProject, err := s.repo.GetByID(ctxT, project.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProjectNotFound
		}
		return err
	}

	if project.Title == "" {
		project.Title = existingProject.Title
	}

	if project.Description == nil {
		project.Description = existingProject.Description
	}

	project.ID = existingProject.ID
	project.OwnerID = existingProject.OwnerID
	project.CreatedAt = existingProject.CreatedAt

	// 2. Обновить проект
	if err := s.repo.Update(ctxT, existingProject); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProjectNotFound
		}
		return err
	}

	return nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 1. Удалить проект
	if err := s.repo.Delete(ctxT, id); err != nil {
		return err
	}

	// 2. Опубликовать событие ProjectDeleted
	event := events.ProjectEvent{
		Type:      events.ProjectDeleted,
		ProjectID: id,
	}
	s.eventBus.Publish(event.ToEvent())

	return nil
}
