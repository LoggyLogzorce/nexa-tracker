package comment

import "gorm.io/gorm"

type Repository interface {
	Create(comment *Comment) error
	GetByID(id uint) (*Comment, error)
	GetByTaskID(taskID uint) ([]Comment, error)
	Update(comment *Comment) error
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(comment *Comment) error {
	// TODO: Implement
	return nil
}

func (r *repository) GetByID(id uint) (*Comment, error) {
	// TODO: Implement
	return nil, nil
}

func (r *repository) GetByTaskID(taskID uint) ([]Comment, error) {
	// TODO: Implement
	return nil, nil
}

func (r *repository) Update(comment *Comment) error {
	// TODO: Implement
	return nil
}

func (r *repository) Delete(id uint) error {
	// TODO: Implement
	return nil
}
