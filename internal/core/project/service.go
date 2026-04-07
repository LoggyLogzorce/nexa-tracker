package project

import "github.com/google/uuid"

type Service interface {
	Create(project *Project, ownerID uuid.UUID) error
	GetByID(id uint) (*Project, error)
	List(ownerID uuid.UUID) ([]Project, error)
	Update(project *Project) error
	Delete(id uint) error
}

type service struct {
	repo            Repository
	statusRepo      interface{} // TODO: Add status repository
	participantRepo interface{} // TODO: Add participant repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(project *Project, ownerID uuid.UUID) error {
	// TODO: Implement
	// 1. Create project
	// 2. Create default statuses (Todo, In Progress, Done)
	// 3. Add owner to participants with role "owner"
	return nil
}

func (s *service) GetByID(id uint) (*Project, error) {
	// TODO: Implement
	return nil, nil
}

func (s *service) List(ownerID uuid.UUID) ([]Project, error) {
	// TODO: Implement
	return nil, nil
}

func (s *service) Update(project *Project) error {
	// TODO: Implement
	return nil
}

func (s *service) Delete(id uint) error {
	// TODO: Implement
	return nil
}
