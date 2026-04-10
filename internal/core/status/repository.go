package status

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(status *Status) error
	CreateBatch(statuses []Status) error
	GetByID(id uint) (*Status, error)
	GetByProjectID(projectID uuid.UUID) ([]Status, error)
	Update(status *Status) error
	Delete(id uint) error
	DeleteByProjectID(projectID uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(status *Status) error {
	// TODO: Implement
	return nil
}

func (r *repository) CreateBatch(statuses []Status) error {
	return r.db.Create(&statuses).Error
}

func (r *repository) GetByID(id uint) (*Status, error) {
	// TODO: Implement
	return nil, nil
}

func (r *repository) GetByProjectID(projectID uuid.UUID) ([]Status, error) {
	// TODO: Implement
	return nil, nil
}

func (r *repository) Update(status *Status) error {
	// TODO: Implement
	return nil
}

func (r *repository) Delete(id uint) error {
	// TODO: Implement
	return nil
}

func (r *repository) DeleteByProjectID(projectID uuid.UUID) error {
	return r.db.Where("project_id = ?", projectID).Delete(&Status{}).Error
}
