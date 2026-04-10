package priority

import "github.com/google/uuid"

type Service interface {
	Create(priority *Priority) error
	GetByID(id uint) (*Priority, error)
	GetByProjectID(projectID uuid.UUID) ([]Priority, error)
	Update(priority *Priority) error
	Delete(id uint) error
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
