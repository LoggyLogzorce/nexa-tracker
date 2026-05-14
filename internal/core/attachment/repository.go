package attachment

import (
	"context"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, attachment *Attachment) error
	GetByID(ctx context.Context, id uint) (*Attachment, error)
	GetByTaskID(ctx context.Context, taskID uint) ([]Attachment, error)
	Delete(ctx context.Context, id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, attachment *Attachment) error {
	return r.db.WithContext(ctx).Create(attachment).Error
}

func (r *repository) GetByID(ctx context.Context, id uint) (*Attachment, error) {
	var a Attachment
	err := r.db.WithContext(ctx).First(&a, id).Error
	return &a, err
}

func (r *repository) GetByTaskID(ctx context.Context, taskID uint) ([]Attachment, error) {
	var attachments []Attachment
	err := r.db.WithContext(ctx).Where("task_id = ?", taskID).Find(&attachments).Error
	return attachments, err
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Attachment{}, id).Error
}
