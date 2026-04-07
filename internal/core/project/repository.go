package project

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(project *Project) error
	GetByID(id uint) (*Project, error)
	List(ownerID uuid.UUID) ([]Project, error)
	Update(project *Project) error
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(project *Project) error {
	// TODO: Implement
	return nil
}

func (r *repository) GetByID(id uint) (*Project, error) {
	// TODO: Implement
	return nil, nil
}

func (r *repository) List(ownerID uuid.UUID) ([]Project, error) {
	// TODO: Implement
	return nil, nil
}

func (r *repository) Update(project *Project) error {
	// TODO: Implement
	return nil
}

func (r *repository) Delete(id uint) error {
	// TODO: Implement
	return nil
}
