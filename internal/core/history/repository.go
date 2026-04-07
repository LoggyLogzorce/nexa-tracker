package history

import "gorm.io/gorm"

type Repository interface {
	Create(history *UpdateHistory) error
	GetByTaskID(taskID uint) ([]UpdateHistory, error)
	GetByUserID(userID string) ([]UpdateHistory, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(history *UpdateHistory) error {
	// TODO: Implement
	return nil
}

func (r *repository) GetByTaskID(taskID uint) ([]UpdateHistory, error) {
	// TODO: Implement
	return nil, nil
}

func (r *repository) GetByUserID(userID string) ([]UpdateHistory, error) {
	// TODO: Implement
	return nil, nil
}
