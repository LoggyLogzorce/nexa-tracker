package project

import (
	"context"
	"github.com/google/uuid"
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
	repo            Repository
	statusRepo      interface{} // TODO: Add status repository
	participantRepo interface{} // TODO: Add participant repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, project *Project, ownerID uuid.UUID) error {
	// TODO: Implement
	// 1. Create project
	// 2. Create default statuses (Todo, In Progress, Done)
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
