package project

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"nexa-task-tracker/internal/core/participant"
	"nexa-task-tracker/internal/core/priority"
	"nexa-task-tracker/internal/core/status"
	"nexa-task-tracker/internal/core/user"
	"nexa-task-tracker/internal/pkg/events"
	"time"
)

type Service interface {
	Create(ctx context.Context, project *Project, ownerID uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*ProjectResponse, error)
	List(ctx context.Context, userID uuid.UUID) ([]ProjectResponse, error)
	ListOwned(ctx context.Context, userID uuid.UUID) ([]ProjectResponse, error)
	Update(ctx context.Context, project UpdateProjectRequest, userID, projectID uuid.UUID) (*Project, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type service struct {
	repo            Repository
	eventBus        *events.EventBus
	userRepo        user.Repository
	statusRepo      status.Repository
	priorityRepo    priority.Repository
	participantRepo participant.Repository
}

func NewService(repo Repository, eventBus *events.EventBus, userRepo user.Repository, statusRepo status.Repository, priorityRepo priority.Repository, participantRepo participant.Repository) Service {
	return &service{
		repo:            repo,
		eventBus:        eventBus,
		userRepo:        userRepo,
		statusRepo:      statusRepo,
		priorityRepo:    priorityRepo,
		participantRepo: participantRepo,
	}
}

var statuses = map[string]bool{
	"plan":        true,
	"in_progress": true,
	"done":        true,
}

var priorities = map[string]bool{
	"low":    true,
	"medium": true,
	"high":   true,
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
		Status:      project.Status,
		Priority:    project.Priority,
	}

	if project.OwnerID == userID {
		projectDto.UserRole = "owner"
	} else {
		participantData, err := s.participantRepo.GetByProjectAndUser(ctxT, project.ID, userID)
		if err == nil {
			projectDto.UserRole = participantData.Role
		}
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

func (s *service) List(ctx context.Context, userID uuid.UUID) ([]ProjectResponse, error) {
	ctxT, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	participants, err := s.participantRepo.GetByUserID(ctxT, userID)
	if err != nil {
		return nil, err
	}

	projectIDs := make([]uuid.UUID, len(participants))
	roleMap := make(map[uuid.UUID]string, len(participants))
	for i, p := range participants {
		projectIDs[i] = p.ProjectID
		roleMap[p.ProjectID] = p.Role
	}

	projects, err := s.repo.ListByParticipant(ctxT, userID, projectIDs)
	if err != nil {
		return nil, err
	}

	ownerIDs := make([]uuid.UUID, 0, len(projects))
	for i := range projects {
		if projects[i].OwnerID == userID {
			continue
		}
		ownerIDs = append(ownerIDs, projects[i].OwnerID)
	}

	ownerMap := make(map[uuid.UUID]user.User)
	if len(ownerIDs) > 0 {
		owners, err := s.userRepo.GetListByIDs(ctxT, ownerIDs)
		if err == nil {
			for i := range owners {
				ownerMap[owners[i].ID] = owners[i]
			}
		}
	}

	currentUser, err := s.userRepo.GetByID(ctxT, userID)
	if err != nil {
		currentUser = nil
	}

	responses := make([]ProjectResponse, 0, len(projects))
	for _, p := range projects {
		resp := ProjectResponse{
			ID:          p.ID,
			Title:       p.Title,
			Description: p.Description,
			CreatedAt:   p.CreatedAt,
			Status:      p.Status,
			Priority:    p.Priority,
			Owner: struct {
				ID    uuid.UUID `json:"id"`
				Name  string    `json:"name"`
				Email string    `json:"email"`
			}{
				ID: p.OwnerID,
			},
		}

		if p.OwnerID == userID && currentUser != nil {
			resp.Owner.Name = currentUser.Name
			resp.Owner.Email = currentUser.Email
			resp.UserRole = "owner"
		} else if p.OwnerID == userID {
			resp.UserRole = "owner"
		} else if owner, ok := ownerMap[p.OwnerID]; ok {
			resp.Owner.Name = owner.Name
			resp.Owner.Email = owner.Email
			resp.UserRole = roleMap[p.ID]
		}

		responses = append(responses, resp)
	}

	return responses, nil
}

func (s *service) ListOwned(ctx context.Context, userID uuid.UUID) ([]ProjectResponse, error) {
	ctxT, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	projects, err := s.repo.ListByOwner(ctxT, userID)
	if err != nil {
		return nil, err
	}

	currentUser, err := s.userRepo.GetByID(ctxT, userID)
	if err != nil {
		currentUser = nil
	}

	responses := make([]ProjectResponse, 0, len(projects))
	for _, p := range projects {
		resp := ProjectResponse{
			ID:          p.ID,
			Title:       p.Title,
			Description: p.Description,
			CreatedAt:   p.CreatedAt,
			Status:      p.Status,
			Priority:    p.Priority,
			UserRole:    "owner",
			Owner: struct {
				ID    uuid.UUID `json:"id"`
				Name  string    `json:"name"`
				Email string    `json:"email"`
			}{
				ID: p.OwnerID,
			},
		}

		if currentUser != nil {
			resp.Owner.Name = currentUser.Name
			resp.Owner.Email = currentUser.Email
		}

		responses = append(responses, resp)
	}

	return responses, nil
}

func (s *service) Update(ctx context.Context, req UpdateProjectRequest, userID, projectID uuid.UUID) (*Project, error) {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 1. Получить проект из БД
	existingProject, err := s.repo.GetByID(ctxT, projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	if req.Title != "" {
		existingProject.Title = req.Title
	}

	if req.Description.Set {
		existingProject.Description = req.Description.Value
	}

	if req.Status.Set {
		if req.Status.Value != nil {
			if _, exists := statuses[*req.Status.Value]; !exists {
				return nil, ErrInvalidStatus
			}
		}
		existingProject.Status = req.Status.Value
	}

	if req.Priority.Set {
		if req.Priority.Value != nil {
			if _, exists := priorities[*req.Priority.Value]; !exists {
				return nil, ErrInvalidPriority
			}
		}
		existingProject.Priority = req.Priority.Value
	}

	// 2. Обновить проект
	if err := s.repo.Update(ctxT, existingProject); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	return existingProject, nil
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
