package status

import "gorm.io/gorm"

type Repository interface {
	Create(status *Status) error
	CreateBatch(statuses []Status) error
	GetByID(id uint) (*Status, error)
	GetByProjectID(projectID uint) ([]Status, error)
	Update(status *Status) error
	Delete(id uint) error
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
	// TODO: Implement
	return nil
}

func (r *repository) GetByID(id uint) (*Status, error) {
	// TODO: Implement
	return nil, nil
}

func (r *repository) GetByProjectID(projectID uint) ([]Status, error) {
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
