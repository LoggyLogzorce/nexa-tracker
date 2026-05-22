package comment

import (
	"context"
	"gorm.io/gorm"
	"nexa-task-tracker/internal/models"
)

type Repository interface {
	Create(ctx context.Context, comment *models.Comment) error
	GetByID(ctx context.Context, id uint) (*models.Comment, error)
	GetByTaskID(ctx context.Context, taskID uint) ([]models.Comment, error)
	Update(ctx context.Context, comment *models.Comment) error
	Delete(ctx context.Context, id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, comment *models.Comment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

func (r *repository) GetByID(ctx context.Context, id uint) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.WithContext(ctx).First(&comment, id).Error
	return &comment, err
}

func (r *repository) GetByTaskID(ctx context.Context, taskID uint) ([]models.Comment, error) {
	var comments []models.Comment
	err := r.db.WithContext(ctx).Find(&comments, "task_id = ?", taskID).Error
	return comments, err
}

func (r *repository) Update(ctx context.Context, comment *models.Comment) error {
	return r.db.WithContext(ctx).Save(comment).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Comment{}, id).Error
}
