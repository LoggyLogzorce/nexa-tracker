package priority

import "gorm.io/gorm"

type Repository interface {
	Create(priority *Priority) error
	GetByID(id uint) (*Priority, error)
	GetByProjectID(projectID uint) ([]Priority, error)
	Update(priority *Priority) error
	Delete(id uint) error
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

func (r *repository) GetByID(id uint) (*Priority, error) {
	// TODO: Implement
	return nil, nil
}

func (r *repository) GetByProjectID(projectID uint) ([]Priority, error) {
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
