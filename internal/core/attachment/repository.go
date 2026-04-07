package attachment

import "gorm.io/gorm"

type Repository interface {
	Create(attachment *Attachment) error
	GetByID(id uint) (*Attachment, error)
	GetByTaskID(taskID uint) ([]Attachment, error)
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(attachment *Attachment) error {
	// TODO: Implement
	return nil
}

func (r *repository) GetByID(id uint) (*Attachment, error) {
	// TODO: Implement
	return nil, nil
}

func (r *repository) GetByTaskID(taskID uint) ([]Attachment, error) {
	// TODO: Implement
	return nil, nil
}

func (r *repository) Delete(id uint) error {
	// TODO: Implement
	return nil
}
