package status

import "github.com/google/uuid"

type Service interface {
	Create(status *Status) error
	GetByID(id uint) (*Status, error)
	GetByProjectID(projectID uuid.UUID) ([]Status, error)
	Update(status *Status) error
	Delete(id uint) error
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
