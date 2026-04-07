package user

import "github.com/google/uuid"

type Service interface {
	Register(email, password, name string) (*User, error)
	GetByID(id uuid.UUID) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Register(email, password, name string) (*User, error) {
	// TODO: Hash password and create user with UUID
	return nil, nil
}

func (s *service) GetByID(id uuid.UUID) (*User, error) {
	// TODO: Implement
	return nil, nil
}

func (s *service) GetByEmail(email string) (*User, error) {
	// TODO: Implement
	return nil, nil
}

func (s *service) Update(user *User) error {
	// TODO: Implement
	return nil
}

func (s *service) Delete(id uuid.UUID) error {
	// TODO: Implement
	return nil
}
