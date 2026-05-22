package attachment

import (
	"context"
	"gorm.io/gorm"
	"nexa-task-tracker/internal/models"
)

type Repository interface {
	Create(ctx context.Context, attachment *models.Attachment) error
	GetByID(ctx context.Context, id uint) (*models.Attachment, error)
	GetByTaskID(ctx context.Context, taskID uint) ([]models.Attachment, error)
	GetByTaskIDs(ctx context.Context, taskIDs []uint) ([]models.Attachment, error)
	Delete(ctx context.Context, id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, attachment *models.Attachment) error {
	return r.db.WithContext(ctx).Create(attachment).Error
}

func (r *repository) GetByID(ctx context.Context, id uint) (*models.Attachment, error) {
	var a models.Attachment
	err := r.db.WithContext(ctx).First(&a, id).Error
	return &a, err
}

func (r *repository) GetByTaskID(ctx context.Context, taskID uint) ([]models.Attachment, error) {
	var attachments []models.Attachment
	err := r.db.WithContext(ctx).Where("task_id = ?", taskID).Find(&attachments).Error
	return attachments, err
}

func (r *repository) GetByTaskIDs(ctx context.Context, taskIDs []uint) ([]models.Attachment, error) {
	var attachments []models.Attachment
	err := r.db.WithContext(ctx).Where("task_id IN ?", taskIDs).Find(&attachments).Error
	return attachments, err
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Attachment{}, id).Error
}
