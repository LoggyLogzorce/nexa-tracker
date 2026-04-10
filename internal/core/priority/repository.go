package priority

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(priority *Priority) error
	CreateBatch(priorities []Priority) error
	GetByID(id uint) (*Priority, error)
	GetByProjectID(projectID uuid.UUID) ([]Priority, error)
	Update(priority *Priority) error
	Delete(id uint) error
	DeleteByProjectID(projectID uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(priority *Priority) error {
	// TODO: Implement
	return nil
}

func (r *repository) CreateBatch(priorities []Priority) error {
	return r.db.Create(&priorities).Error
}

func (r *repository) GetByID(id uint) (*Priority, error) {
	// TODO: Implement
	return nil, nil
}

func (r *repository) GetByProjectID(projectID uuid.UUID) ([]Priority, error) {
	// TODO: Implement
	return nil, nil
}

func (r *repository) Update(priority *Priority) error {
	// TODO: Implement
	return nil
}

func (r *repository) Delete(id uint) error {
	// TODO: Implement
	return nil
}

func (r *repository) DeleteByProjectID(projectID uuid.UUID) error {
	return r.db.Where("project_id = ?", projectID).Delete(&Priority{}).Error
}
